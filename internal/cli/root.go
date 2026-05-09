package cli

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"vigtem/internal/config"
	"vigtem/internal/debugger"
	"vigtem/internal/engine"
)

func Execute() {
	target := flag.String("target", "", "carpeta a analizar")
	short := flag.String("t", "", "carpeta a analizar")
	flag.Parse()

	chosen := *target
	if chosen == "" {
		chosen = *short
	}
	if chosen == "" {
		fmt.Print("📁 Escribe la carpeta a examinar: ")
		in := bufio.NewScanner(os.Stdin)
		if !in.Scan() {
			fmt.Fprintln(os.Stderr, "no se recibio carpeta")
			os.Exit(1)
		}
		chosen = in.Text()
	}

	chosen = filepath.Clean(chosen)
	if st, err := os.Stat(chosen); err != nil || !st.IsDir() {
		fmt.Fprintf(os.Stderr, "ruta invalida: %s\n", chosen)
		os.Exit(1)
	}

	cfg := config.Default(chosen)
	tracer := debugger.New(os.Stdout)
	eng := engine.New(cfg, tracer)
	rep, err := eng.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ Analisis finalizado\nWorkspace: %s\nArchivos inspeccionados: %d\nEventos de traza: %d\n", rep.Workspace, len(rep.Findings), rep.TraceEvents)
}
