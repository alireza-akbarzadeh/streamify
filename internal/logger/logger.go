package logger

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// ANSI color codes
const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorGray    = "\033[90m"
	colorWhite   = "\033[97m"
	colorBold    = "\033[1m"
)

// ContextKey is the type for context keys
type ContextKey string

const (
	// RequestIDKey is the context key for request ID
	RequestIDKey ContextKey = "request_id"
	// UserIDKey is the context key for user ID (string, to match middleware)
	UserIDKey ContextKey = "user_id"
)

var (
	// Logger is the global logger instance
	Logger *slog.Logger
)

// init automatically initializes the logger when the package is imported
func init() {
	InitLogger()
}

// PrettyHandler is a custom slog handler for pretty console output
type PrettyHandler struct {
	opts   slog.HandlerOptions
	mu     *sync.Mutex
	out    io.Writer
	attrs  []slog.Attr
	groups []string
}

// NewPrettyHandler creates a new pretty handler
func NewPrettyHandler(out io.Writer, opts *slog.HandlerOptions) *PrettyHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &PrettyHandler{
		opts: *opts,
		mu:   &sync.Mutex{},
		out:  out,
	}
}

func (h *PrettyHandler) Enabled(_ context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Format timestamp
	timeStr := r.Time.Format("15:04:05")

	// Get level with color and emoji
	levelStr, levelColor := h.formatLevel(r.Level)

	// Build the log line
	var sb strings.Builder

	// Time in gray
	sb.WriteString(colorGray)
	sb.WriteString(timeStr)
	sb.WriteString(colorReset)
	sb.WriteString(" ")

	// Level with color
	sb.WriteString(levelColor)
	sb.WriteString(levelStr)
	sb.WriteString(colorReset)
	sb.WriteString(" ")

	// Message in white/bold
	sb.WriteString(colorWhite)
	sb.WriteString(r.Message)
	sb.WriteString(colorReset)

	// Attributes
	if r.NumAttrs() > 0 || len(h.attrs) > 0 {
		sb.WriteString(" ")
		sb.WriteString(colorGray)

		// Pre-stored attributes
		for _, attr := range h.attrs {
			h.writeAttr(&sb, attr)
		}

		// Record attributes
		r.Attrs(func(a slog.Attr) bool {
			h.writeAttr(&sb, a)
			return true
		})

		sb.WriteString(colorReset)
	}

	sb.WriteString("\n")

	_, err := h.out.Write([]byte(sb.String()))
	return err
}

func (h *PrettyHandler) writeAttr(sb *strings.Builder, a slog.Attr) {
	if a.Equal(slog.Attr{}) {
		return
	}

	sb.WriteString(colorCyan)
	sb.WriteString(a.Key)
	sb.WriteString(colorReset)
	sb.WriteString("=")
	sb.WriteString(colorYellow)
	sb.WriteString(fmt.Sprintf("%v", a.Value.Any()))
	sb.WriteString(colorReset)
	sb.WriteString(" ")
}

func (h *PrettyHandler) formatLevel(level slog.Level) (string, string) {
	switch {
	case level >= slog.LevelError:
		return "❌ ERROR", colorRed
	case level >= slog.LevelWarn:
		return "⚠️  WARN ", colorYellow
	case level >= slog.LevelInfo:
		return "ℹ️  INFO ", colorBlue
	default:
		return "🔍 DEBUG", colorMagenta
	}
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	copy(newAttrs[len(h.attrs):], attrs)
	return &PrettyHandler{
		opts:   h.opts,
		mu:     h.mu,
		out:    h.out,
		attrs:  newAttrs,
		groups: h.groups,
	}
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	newGroups := make([]string, len(h.groups)+1)
	copy(newGroups, h.groups)
	newGroups[len(h.groups)] = name
	return &PrettyHandler{
		opts:   h.opts,
		mu:     h.mu,
		out:    h.out,
		attrs:  h.attrs,
		groups: newGroups,
	}
}

// InitLogger initializes the global logger
func InitLogger() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	var handler slog.Handler

	env := os.Getenv("GO_ENV")
	if env == "development" || env == "dev" || env == "local" {
		opts.Level = slog.LevelDebug
		// Use pretty handler in development
		handler = NewPrettyHandler(os.Stdout, opts)
	} else {
		// Use JSON handler in production
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	Logger = slog.New(handler)
	slog.SetDefault(Logger)
}

// generateRequestID generates a simple request ID
func generateRequestID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// WithRequestID adds a request ID to the context
func WithRequestID(ctx context.Context) context.Context {
	requestID := generateRequestID()
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// WithUserID adds a user ID (string) to the context
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// GetRequestID gets the request ID from context
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}

// GetUserID gets the user ID (string) from context
func GetUserID(ctx context.Context) string {
	if id, ok := ctx.Value(UserIDKey).(string); ok {
		return id
	}
	return ""
}

// LogWithContext creates a logger with context values
func LogWithContext(ctx context.Context) *slog.Logger {
	logger := Logger

	// If Logger is not initialized, use the default slog logger
	if logger == nil {
		logger = slog.Default()
	}

	if requestID := GetRequestID(ctx); requestID != "" {
		logger = logger.With("request_id", requestID)
	}

	if userID := GetUserID(ctx); userID != "" {
		logger = logger.With("user_id", userID)
	}

	return logger
}

// Info logs an info message with context
func Info(ctx context.Context, msg string, args ...any) {
	LogWithContext(ctx).Info(msg, args...)
}

// Error logs an error message with context
func Error(ctx context.Context, msg string, err error, args ...any) {
	logger := LogWithContext(ctx)
	if err != nil {
		logger = logger.With("error", err.Error())
	}
	logger.Error(msg, args...)
}

// Debug logs a debug message with context
func Debug(ctx context.Context, msg string, args ...any) {
	LogWithContext(ctx).Debug(msg, args...)
}

// Warn logs a warning message with context
func Warn(ctx context.Context, msg string, args ...any) {
	LogWithContext(ctx).Warn(msg, args...)
}

// RequestLoggerMiddleware is a chi-compatible middleware that adds request ID to context and logs requests
func RequestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Add request ID to context
		ctx := WithRequestID(r.Context())
		r = r.WithContext(ctx)

		// Wrap the ResponseWriter to capture status code
		rw := &responseWriter{ResponseWriter: w, status: 200}

		// Process request
		next.ServeHTTP(rw, r)

		// Log request completion
		duration := time.Since(start)
		LogWithContext(ctx).Info("request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.status,
			"duration_ms", duration.Milliseconds(),
			"ip", r.RemoteAddr,
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
