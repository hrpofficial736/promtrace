package store

import (
	"database/sql"
	"time"

	"github.com/hrpofficial736/promtrace/internal/logger"
	_ "github.com/mattn/go-sqlite3"
)

type Store interface {
	SaveTrace(trace *Trace) error
	GetTrace(id string) (*Trace, error)
	ListTraces(limit int) ([]*Trace, error)
	ListAllTraces() ([]*Trace, error)
	ListSessions() ([]*Session, error)
	GetStats(since time.Duration) (*Stats, error)
	Close() error
}

type sqliteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dbPath string) (Store, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS traces (
		id 			  TEXT	PRIMARY KEY, 
		session_id 	  TEXT,
		timestamp     DATETIME,
		host          TEXT,
		method        TEXT,
		path          TEXT,
		model		  TEXT,
		tokens 		  INTEGER,
		cost 		  INTEGER,
		system_prompt TEXT,
		user_prompt   TEXT,
		request_body  TEXT,
		response      TEXT,
		status_code   INTEGER,
		latency_ms    INTEGER,
		created_at    DATETIME
	)
	`)

	if err != nil {
		return nil, err
	}

	return &sqliteStore{db: db}, nil

}

func (s *sqliteStore) SaveTrace(trace *Trace) error {
	_, err := s.db.Exec(
		`INSERT INTO traces (id, session_id, timestamp, host, method, path, model, tokens, cost, system_prompt, user_prompt, request_body, response, status_code, latency_ms, created_at)
         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		trace.ID,
		trace.SessionID,
		trace.Timestamp,
		trace.Host,
		trace.Method,
		trace.Path,
		trace.Model,
		trace.Tokens,
		trace.Cost,
		trace.SystemPrompt,
		trace.UserPrompt,
		trace.RequestBody,
		trace.Response,
		trace.StatusCode,
		trace.LatencyMs,
		trace.CreatedAt,
	)

	return err
}

func (s *sqliteStore) GetTrace(id string) (*Trace, error) {
	var trace *Trace = &Trace{}
	row := s.db.QueryRow(`
	SELECT * FROM traces WHERE id = ?
	`, id)

	err := row.Scan(
		&trace.ID,
		&trace.SessionID,
		&trace.Timestamp,
		&trace.Host,
		&trace.Method,
		&trace.Path,
		&trace.Model,
		&trace.Tokens,
		&trace.Cost,
		&trace.SystemPrompt,
		&trace.UserPrompt,
		&trace.RequestBody,
		&trace.Response,
		&trace.StatusCode,
		&trace.LatencyMs,
		&trace.CreatedAt,
	)

	if err != nil {
		logger.Log.Error("error while scanning rows in getting traces", "error", err)
		return nil, err
	}

	return trace, nil
}

func (s *sqliteStore) ListTraces(limit int) ([]*Trace, error) {
	var traces []*Trace

	rows, err := s.db.Query(`
	SELECT * FROM traces ORDER BY timestamp DESC LIMIT ?
	`, limit)

	if err != nil {
		logger.Log.Error("error while getting all the traces in listing traces", "error", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var trace *Trace = &Trace{}
		err := rows.Scan(
			&trace.ID,
			&trace.SessionID,
			&trace.Timestamp,
			&trace.Host,
			&trace.Method,
			&trace.Path,
			&trace.Model,
			&trace.Tokens,
			&trace.Cost,
			&trace.SystemPrompt,
			&trace.UserPrompt,
			&trace.RequestBody,
			&trace.Response,
			&trace.StatusCode,
			&trace.LatencyMs,
			&trace.CreatedAt,
		)
		if err != nil {
			logger.Log.Error("error while scanning rows in listing traces", "error", err)
		}
		traces = append(traces, trace)
	}

	return traces, nil
}

func (s *sqliteStore) ListAllTraces() ([]*Trace, error) {
	var traces []*Trace

	rows, err := s.db.Query(`SELECT * FROM traces ORDER BY timestamp ASC`)

	if err != nil {
		logger.Log.Error("error while getting all the traces in listing traces", "error", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var trace *Trace = &Trace{}
		err := rows.Scan(
			&trace.ID,
			&trace.SessionID,
			&trace.Timestamp,
			&trace.Host,
			&trace.Method,
			&trace.Path,
			&trace.Model,
			&trace.Tokens,
			&trace.Cost,
			&trace.SystemPrompt,
			&trace.UserPrompt,
			&trace.RequestBody,
			&trace.Response,
			&trace.StatusCode,
			&trace.LatencyMs,
			&trace.CreatedAt,
		)
		if err != nil {
			logger.Log.Error("error while scanning rows in listing traces", "error", err)
		}
		traces = append(traces, trace)
	}

	return traces, nil

}

func (s *sqliteStore) ListSessions() ([]*Session, error) {
	var sessions []*Session

	rows, err := s.db.Query(`
	SELECT session_id, COUNT(*) as total_calls,
	MIN(timestamp) as started_at, AVG(latency_ms) as avg_latency,
	SUM(tokens) as total_tokens, SUM(cost) as total_cost
	FROM traces GROUP BY session_id ORDER BY started_at DESC
	`)

	if err != nil {
		logger.Log.Error("error while fetching sessions from DB", "error", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var session *Session = &Session{}

		sqliteLayout := "2006-01-02 15:04:05.999999999-07:00"
		var s_at_string string

		err := rows.Scan(&session.ID, &session.TotalCalls, &s_at_string, &session.AvgLatency, &session.TotalTokens, &session.TotalCost)

		if err != nil {
			logger.Log.Error("error while scanning sessions from rows", "error", err)
			return nil, err
		}

		logger.Log.Info(s_at_string)

		session.StartedAt, err = time.Parse(sqliteLayout, s_at_string)

		if err != nil {
			logger.Log.Error("error while parsing string to time type", "error", err)
			return nil, err
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (s *sqliteStore) GetStats(since time.Duration) (*Stats, error) {
	stats := &Stats{}

	cutoff := time.Now().Add(-since)

	err := s.db.QueryRow(`
	SELECT COUNT(*), SUM(tokens),
	SUM(cost), AVG(latency_ms) FROM traces
	WHERE timestamp >= ?
	`, cutoff).Scan(&stats.TotalCalls, &stats.TotalTokens, &stats.TotalCost, &stats.AvgLatency)

	if err != nil {
		logger.Log.Error("error while fetching stats from db", "error", err)
		return nil, err
	}

	rows, err := s.db.Query(`
	SELECT DATE(timestamp) as day, COUNT(*),
	SUM(tokens), SUM(cost), AVG(latency_ms)
	FROM traces WHERE timestamp >= ?
	GROUP BY DATE(timestamp) ORDER BY day ASC
	`, cutoff)

	if err != nil {
		logger.Log.Error("error while fetching stats from db", "error", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var t trendData
		rows.Scan(&t.Date, &t.Calls, &t.Tokens, &t.Cost, &t.AvgLatency)

		stats.Trend = append(stats.Trend, &t)
	}

	return stats, nil
}

func (s *sqliteStore) Close() error {
	return s.db.Close()
}
