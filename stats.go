package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"runtime"

	"golang.zx2c4.com/wireguard/device"
	"gopkg.in/ini.v1"
)

func ToJson[T any](i *T) string {
	s, _ := json.MarshalIndent(i, "", "  ")
	return string(s)
}

func Stats(w http.ResponseWriter, dev *device.Device) {
	var buf bytes.Buffer
	if err := dev.IpcGetOperation(&buf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stats := struct {
		Endpoint               string `ini:"endpoint"`
		LastHandshakeTimestamp int64  `ini:"last_handshake_time_sec"`
		ReceivedBytes          int64  `ini:"rx_bytes"`
		SentBytes              int64  `ini:"tx_bytes"`
		NumGoroutine           int
	}{NumGoroutine: runtime.NumGoroutine()}

	ini.MapTo(&stats, buf.Bytes())

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, ToJson(&stats))
}
