package http

import (
	"context"
	"errors"
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
	httpServer := &http.Server{
		Addr: staticConfig.GetPortString(),
	}
	handler := NewHandler(mcpServer, httpServer, staticConfig.SSEBaseURL)
	httpServer.Handler = RequestMiddleware(handler)

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
	defer stop()

	serverErr := make(chan error, 1)
	go func() {
		log.Info().
			Str("addr", httpServer.Addr).
			Str("mcp", MCPEndpoint).
			Str("sse", SSEEndpoint).
			Str("message", SSEMessageEndpoint).
			Msg("starting HTTP MCP server")
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case <-ctx.Done():
	case err := <-serverErr:
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := handler.Shutdown(shutdownCtx); err != nil {
		return err
	}
	return nil
}

// Handler owns mounted MCP HTTP transports so they can be shut down cleanly.
type Handler struct {
	mux                  *http.ServeMux
	shutdownCtx          context.Context
	shutdownCancel       context.CancelFunc
	sseServer            *mcpgo.SSEServer
	streamableHTTPServer *mcpgo.StreamableHTTPServer
}

// NewHandler wires HTTP routes to MCP transport handlers.
func NewHandler(mcpServer *mcpserver.Server, httpServer *http.Server, sseBaseURL string) *Handler {
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
	}

	mux.Handle(SSEEndpoint, sseServer.SSEHandler())
	mux.Handle(SSEMessageEndpoint, sseServer.MessageHandler())
	mux.Handle(MCPEndpoint, handler.withShutdownContext(streamableHTTPServer))
	mux.HandleFunc(HealthEndpoint, func(w http.ResponseWriter, r *http.Request) {
		if mcpServer.IsHealthy() {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("healthy"))
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("unhealthy"))
	})

	return handler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) withShutdownContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(r.Context())
		done := make(chan struct{})
		go func() {
			select {
			case <-h.shutdownCtx.Done():
				cancel()
			case <-done:
			}
		}()
		defer func() {
			close(done)
			cancel()
		}()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Shutdown closes active MCP transport sessions and the underlying HTTP server.
func (h *Handler) Shutdown(ctx context.Context) error {
	h.shutdownCancel()

	var shutdownErr error
	if h.sseServer != nil {
		shutdownErr = h.sseServer.Shutdown(ctx)
	}
	if h.streamableHTTPServer != nil {
		if err := h.streamableHTTPServer.Shutdown(ctx); shutdownErr == nil {
			shutdownErr = err
		}
	}
	return shutdownErr
}
