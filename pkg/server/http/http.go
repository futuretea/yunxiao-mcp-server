package http

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	mcpgo "github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"

	"github.com/futuretea/yunxiao-mcp-server/pkg/core/config"
	mcpserver "github.com/futuretea/yunxiao-mcp-server/pkg/server/mcp"
)

const (
	HealthEndpoint     = "/healthz"
	MCPEndpoint        = "/mcp"
	SSEEndpoint        = "/sse"
	SSEMessageEndpoint = "/message"
)

// Serve starts HTTP transports for MCP and blocks until shutdown.
func Serve(ctx context.Context, mcpServer *mcpserver.Server, staticConfig *config.StaticConfig) error {
	if staticConfig == nil {
		return errors.New("static config is required")
	}

	listener, err := net.Listen("tcp", staticConfig.GetPortString())
	if err != nil {
		return err
	}
	defer func() {
		if err := listener.Close(); err != nil {
			log.Error().Err(err).Msg("close listener")
		}
	}()

	return ServeListener(ctx, mcpServer, staticConfig, listener)
}

const (
	defaultReadHeaderTimeout = 5 * time.Second
	defaultReadTimeout       = 30 * time.Second
	defaultWriteTimeout      = 30 * time.Second
	defaultIdleTimeout       = 120 * time.Second
)

// ServeListener starts HTTP transports for MCP on the provided listener and blocks until shutdown.
func ServeListener(ctx context.Context, mcpServer *mcpserver.Server, staticConfig *config.StaticConfig, listener net.Listener) error {
	if err := validateServeListenerInputs(mcpServer, staticConfig, listener); err != nil {
		return err
	}
	httpServer, handler := newServerWithHandler(mcpServer, staticConfig, listener)
	return runServer(ctx, httpServer, handler, listener)
}

func validateServeListenerInputs(mcpServer *mcpserver.Server, staticConfig *config.StaticConfig, listener net.Listener) error {
	if mcpServer == nil {
		return errors.New("MCP server is required")
	}
	if staticConfig == nil {
		return errors.New("static config is required")
	}
	if listener == nil {
		return errors.New("listener is required")
	}
	return nil
}

func newServerWithHandler(mcpServer *mcpserver.Server, staticConfig *config.StaticConfig, listener net.Listener) (*http.Server, *Handler) {
	httpServer := &http.Server{
		Addr:              listener.Addr().String(),
		ReadHeaderTimeout: defaultReadHeaderTimeout,
		ReadTimeout:       defaultReadTimeout,
		WriteTimeout:      defaultWriteTimeout,
		IdleTimeout:       defaultIdleTimeout,
	}
	handler := NewHandler(mcpServer, httpServer, staticConfig.SSEBaseURL)
	httpServer.Handler = RequestMiddleware(handler)
	return httpServer, handler
}

func runServer(ctx context.Context, httpServer *http.Server, handler *Handler, listener net.Listener) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
	defer stop()

	serverErr := make(chan error, 1)
	serverDone := make(chan struct{})
	go func() {
		defer close(serverDone)
		log.Info().
			Str("addr", httpServer.Addr).
			Str("mcp", MCPEndpoint).
			Str("sse", SSEEndpoint).
			Str("message", SSEMessageEndpoint).
			Msg("starting HTTP MCP server")
		if err := httpServer.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	var serveErr error
	select {
	case <-ctx.Done():
	case err := <-serverErr:
		serveErr = err
	}

	select {
	case err := <-serverErr:
		if serveErr == nil {
			serveErr = err
		}
	default:
	}

	if err := handler.Shutdown(9 * time.Second); err != nil {
		if serveErr != nil {
			return errors.Join(serveErr, err)
		}
		return err
	}
	<-serverDone
	return serveErr
}

// Handler owns mounted MCP HTTP transports so they can be shut down cleanly.
type Handler struct {
	mux                  *http.ServeMux
	shutdownCtx          context.Context
	shutdownCancel       context.CancelFunc
	sseServer            *mcpgo.SSEServer
	streamableHTTPServer *mcpgo.StreamableHTTPServer
	httpServer           *http.Server
}

// NewHandler wires HTTP routes to MCP transport handlers.
func NewHandler(mcpServer *mcpserver.Server, httpServer *http.Server, sseBaseURL string) *Handler {
	if mcpServer == nil {
		panic("mcp server is required")
	}
	if httpServer == nil {
		panic("http server is required")
	}

	mux := http.NewServeMux()
	shutdownCtx, shutdownCancel := context.WithCancel(context.Background())

	sseServer := mcpServer.ServeSSE(sseBaseURL, httpServer)
	streamableHTTPServer := mcpServer.ServeStreamableHTTP(httpServer)

	handler := &Handler{
		mux:                  mux,
		shutdownCtx:          shutdownCtx,
		shutdownCancel:       shutdownCancel,
		sseServer:            sseServer,
		streamableHTTPServer: streamableHTTPServer,
		httpServer:           httpServer,
	}

	mux.Handle(SSEEndpoint, sseServer.SSEHandler())
	mux.Handle(SSEMessageEndpoint, sseServer.MessageHandler())
	mux.Handle(MCPEndpoint, handler.withShutdownContext(streamableHTTPServer))
	mux.HandleFunc(HealthEndpoint, func(w http.ResponseWriter, r *http.Request) {
		status, body := http.StatusServiceUnavailable, "unhealthy"
		if mcpServer.IsHealthy() {
			status, body = http.StatusOK, "healthy"
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	})

	return handler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) withShutdownContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()
		go func() {
			select {
			case <-h.shutdownCtx.Done():
				cancel()
			case <-ctx.Done():
			}
		}()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Shutdown closes active MCP transport sessions and the underlying HTTP server.
// Each active subsystem receives an equal slice of the supplied timeout budget.
func (h *Handler) Shutdown(timeout time.Duration) error {
	h.shutdownCancel()

	active := 0
	if h.sseServer != nil {
		active++
	}
	if h.streamableHTTPServer != nil {
		active++
	}
	if h.httpServer != nil {
		active++
	}
	if active == 0 {
		return nil
	}

	const minSlice = 50 * time.Millisecond
	perCall := timeout / time.Duration(active)
	if perCall < minSlice {
		perCall = minSlice
	}

	var errs []error
	if h.sseServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), perCall)
		if err := h.sseServer.Shutdown(ctx); err != nil {
			errs = append(errs, err)
		}
		cancel()
	}
	if h.streamableHTTPServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), perCall)
		if err := h.streamableHTTPServer.Shutdown(ctx); err != nil {
			errs = append(errs, err)
		}
		cancel()
	}
	if h.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), perCall)
		if err := h.httpServer.Shutdown(ctx); err != nil {
			errs = append(errs, err)
		}
		cancel()
	}
	return errors.Join(errs...)
}
