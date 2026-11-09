package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/user/j2cpp-engine/parser"
)

type TranspileRequest struct {
	Source string `json:"source"`
}

type TranspileResponse struct {
	Success bool   `json:"success"`
	Header  string `json:"header"`
	Cpp     string `json:"cpp"`
	Error   string `json:"error,omitempty"`
}

type StatusResponse struct {
	Version        string   `json:"version"`
	RuntimeHeaders []string `json:"runtime_headers"`
	Features       []string `json:"features"`
}

const dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>libjruntime - Transpilation Studio</title>
    <style>
        :root {
            --bg-primary: #1e1e2e;
            --bg-secondary: #252538;
            --accent: #89b4fa;
            --text-main: #cdd6f4;
            --text-muted: #a6adc8;
            --border: #45475a;
            --success: #a6e3a1;
            --danger: #f38ba8;
        }
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            background: var(--bg-primary);
            color: var(--text-main);
            height: 100vh;
            display: flex;
            flex-direction: column;
        }
        header {
            background: var(--bg-secondary);
            border-bottom: 1px solid var(--border);
            padding: 12px 24px;
            display: flex;
            align-items: center;
            justify-content: space-between;
        }
        .brand {
            display: flex;
            align-items: center;
            gap: 12px;
            font-weight: 700;
            font-size: 1.15rem;
            color: var(--accent);
        }
        .badge {
            background: rgba(137, 180, 250, 0.2);
            color: var(--accent);
            padding: 2px 8px;
            border-radius: 12px;
            font-size: 0.8rem;
        }
        .header-actions {
            display: flex;
            gap: 12px;
            align-items: center;
        }
        button.btn {
            background: var(--accent);
            color: #11111b;
            border: none;
            padding: 8px 16px;
            border-radius: 6px;
            font-weight: 600;
            cursor: pointer;
            transition: opacity 0.2s;
        }
        button.btn:hover { opacity: 0.9; }
        .container {
            flex: 1;
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 16px;
            padding: 16px 24px;
            min-height: 0;
        }
        .panel {
            background: var(--bg-secondary);
            border: 1px solid var(--border);
            border-radius: 8px;
            display: flex;
            flex-direction: column;
            overflow: hidden;
        }
        .panel-header {
            padding: 10px 16px;
            background: rgba(0,0,0,0.2);
            border-bottom: 1px solid var(--border);
            display: flex;
            justify-content: space-between;
            align-items: center;
            font-size: 0.9rem;
            font-weight: 600;
        }
        .editor-area {
            flex: 1;
            display: flex;
            flex-direction: column;
        }
        textarea {
            flex: 1;
            width: 100%;
            background: transparent;
            color: var(--text-main);
            border: none;
            outline: none;
            padding: 16px;
            font-family: ui-monospace, SFMono-Regular, "SF Mono", Menlo, Consolas, monospace;
            font-size: 0.95rem;
            resize: none;
            line-height: 1.5;
        }
        .tooltip-marker {
            display: inline-flex;
            align-items: center;
            justify-content: center;
            width: 16px;
            height: 16px;
            border-radius: 50%;
            background: var(--border);
            color: var(--text-muted);
            font-size: 0.75rem;
            cursor: help;
            margin-left: 6px;
        }
        .tabs {
            display: flex;
            gap: 4px;
        }
        .tab {
            padding: 4px 12px;
            border-radius: 4px;
            cursor: pointer;
            background: transparent;
            color: var(--text-muted);
            font-size: 0.85rem;
        }
        .tab.active {
            background: var(--border);
            color: var(--text-main);
        }
        footer {
            background: var(--bg-secondary);
            border-top: 1px solid var(--border);
            padding: 8px 24px;
            font-size: 0.8rem;
            color: var(--text-muted);
            display: flex;
            justify-content: space-between;
        }
    </style>
</head>
<body>
    <header>
        <div class="brand">
            <span>⚡ libjruntime</span>
            <span class="badge" id="version-badge">v1.0.38</span>
        </div>
        <div class="header-actions">
            <button class="btn" onclick="transpileCode()">Transpile Code</button>
        </div>
    </header>

    <div class="container">
        <div class="panel">
            <div class="panel-header">
                <div>
                    <span>Java Source Code</span>
                    <span class="tooltip-marker" title="Input standard Java class, interface, method, lambda, or collection usage.">?</span>
                </div>
            </div>
            <div class="editor-area">
                <textarea id="java-input" spellcheck="false">public class ExampleApp {
    public static void main(String[] args) {
        System.out.println("Hello, libjruntime!");
        auto incrementer = (x) -> x + 1;
        int res = incrementer(41);
        System.out.println(res);
    }
}</textarea>
            </div>
        </div>

        <div class="panel">
            <div class="panel-header">
                <div>
                    <span>Transpiled C++20</span>
                    <span class="tooltip-marker" title="C++ code emitted by compiler-engine, linked against libjruntime and Boehm GC.">?</span>
                </div>
                <div class="tabs">
                    <span class="tab active" id="tab-cpp" onclick="switchTab('cpp')">Implementation (.cpp)</span>
                    <span class="tab" id="tab-h" onclick="switchTab('h')">Header (.h)</span>
                </div>
            </div>
            <div class="editor-area">
                <textarea id="output" readonly spellcheck="false">// Click "Transpile Code" to inspect output...</textarea>
            </div>
        </div>
    </div>

    <footer>
        <div id="status-bar">Ready - Memory Management: Boehm GC (libgc)</div>
        <div>Modern C++20 Zero-Cost Runtime</div>
    </footer>

    <script>
        let latestResult = { header: "", cpp: "" };
        let activeTab = 'cpp';

        async function transpileCode() {
            const input = document.getElementById("java-input").value;
            const outputEl = document.getElementById("output");
            const statusEl = document.getElementById("status-bar");
            statusEl.innerText = "Transpiling...";

            try {
                const res = await fetch("/api/transpile", {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({ source: input })
                });
                const data = await res.json();
                if (data.success) {
                    latestResult.header = data.header;
                    latestResult.cpp = data.cpp;
                    renderTab();
                    statusEl.innerText = "Transpilation succeeded";
                } else {
                    outputEl.value = "Error: " + data.error;
                    statusEl.innerText = "Transpilation error";
                }
            } catch (err) {
                outputEl.value = "Network error: " + err;
                statusEl.innerText = "Connection failed";
            }
        }

        function switchTab(tab) {
            activeTab = tab;
            document.getElementById("tab-cpp").classList.toggle("active", tab === 'cpp');
            document.getElementById("tab-h").classList.toggle("active", tab === 'h');
            renderTab();
        }

        function renderTab() {
            const outputEl = document.getElementById("output");
            outputEl.value = activeTab === 'cpp' ? latestResult.cpp : latestResult.header;
        }

        async function fetchStatus() {
            try {
                const res = await fetch("/api/status");
                const data = await res.json();
                document.getElementById("version-badge").innerText = data.version;
            } catch (e) {}
        }
        fetchStatus();
    </script>
</body>
</html>`

func handleTranspileSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TranspileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	tempDir, err := os.MkdirTemp("", "j2cpp_dash_*")
	if err != nil {
		json.NewEncoder(w).Encode(TranspileResponse{Success: false, Error: err.Error()})
		return
	}
	defer os.RemoveAll(tempDir)

	input := antlr.NewInputStream(req.Source)
	lexer := parser.NewJavaLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewJavaParser(stream)
	tree := p.CompilationUnit()

	visitor := NewCppTranspilerVisitor(tempDir)
	visitor.Visit(tree)

	var headerContent string
	var cppContent string

	entries, _ := os.ReadDir(tempDir)
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".h") {
			data, _ := os.ReadFile(filepath.Join(tempDir, entry.Name()))
			headerContent = string(data)
		} else if strings.HasSuffix(entry.Name(), ".cpp") {
			data, _ := os.ReadFile(filepath.Join(tempDir, entry.Name()))
			cppContent = string(data)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TranspileResponse{
		Success: true,
		Header:  headerContent,
		Cpp:     cppContent,
	})
}

func getProjectVersion() string {
	for _, p := range []string{"../VERSION", "VERSION"} {
		if data, err := os.ReadFile(p); err == nil {
			v := strings.TrimSpace(string(data))
			if v != "" {
				return "v" + v
			}
		}
	}
	return "v1.0.62"
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(StatusResponse{
		Version: getProjectVersion(),
		RuntimeHeaders: []string{
			"Object.h", "String.h", "StringBuilder.h", "Array.h", "ArrayList.h",
			"LinkedList.h", "Collections.h", "HashMap.h", "HashSet.h", "Boxes.h", "Math.h", "System.h",
			"Exceptions.h", "FinallyGuard.h", "Thread.h", "Atomic.h", "Lock.h", "Functions.h",
			"Stream.h", "Class.h", "Optional.h",
		},
		Features: []string{
			"Class & Interface Inheritance", "Constructors & Chaining",
			"Exception Handling (try-catch-finally)", "Lambdas & Method References",
			"Streams & Functional Interfaces", "Reflection (Class<T>)",
			"Collections Algorithms (sort, reverse, etc.)", "Boehm GC (libgc)",
		},
	})
}

func startDashboardServer(port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		html := strings.ReplaceAll(dashboardHTML, "v1.0.38", getProjectVersion())
		w.Write([]byte(html))
	})
	mux.HandleFunc("/api/transpile", handleTranspileSnippet)
	mux.HandleFunc("/api/status", handleStatus)

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Starting libjruntime interactive dashboard on http://localhost%s\n", addr)
	return http.ListenAndServe(addr, mux)
}
