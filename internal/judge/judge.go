package judge

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Judge struct {
	goPath      string
	gccPath     string
	gppPath     string
	javacPath   string
	javaPath    string
	timeLimit   int // 毫秒
	memoryLimit int // MB
}

func NewJudge(goPath string, timeLimit, memoryLimit int) *Judge {
	return &Judge{
		goPath:      goPath,
		gccPath:     "/usr/bin/gcc",
		gppPath:     "/usr/bin/g++",
		javacPath:   "/usr/bin/javac",
		javaPath:    "/usr/bin/java",
		timeLimit:   timeLimit,
		memoryLimit: memoryLimit,
	}
}

func NewJudgeWithPaths(goPath, gccPath, gppPath, javacPath, javaPath string, timeLimit, memoryLimit int) *Judge {
	return &Judge{
		goPath:      goPath,
		gccPath:     gccPath,
		gppPath:     gppPath,
		javacPath:   javacPath,
		javaPath:    javaPath,
		timeLimit:   timeLimit,
		memoryLimit: memoryLimit,
	}
}

// JudgeResult 判题结果
type JudgeResult struct {
	Status     string
	Output     string
	Expected   string
	Error      string
	TimeUsed   int // 毫秒
	MemoryUsed int // KB
}

// CompiledProgram 编译后的程序（支持多次执行）
type CompiledProgram struct {
	TmpDir      string // 临时目录路径
	BinaryPath  string // 可执行文件路径（Go/C/C++）
	IsJava      bool   // 是否为 Java 程序
	JavaClassDir string // Java class 文件所在目录
}

// Cleanup 清理临时目录
func (p *CompiledProgram) Cleanup() {
	if p.TmpDir != "" {
		os.RemoveAll(p.TmpDir)
	}
}

// Compile 编译代码，返回可复用的 CompiledProgram
// 调用方负责调用 Cleanup() 清理资源
func (j *Judge) Compile(language, code, driverCode string) (*CompiledProgram, *JudgeResult) {
	switch language {
	case "C":
		return j.compileC(code)
	case "C++":
		return j.compileCpp(code)
	case "Java":
		return j.compileJava(code)
	default:
		return j.compileGo(code, driverCode)
	}
}

// RunCompiled 执行已编译的程序（一次编译，多次执行）
func (j *Judge) RunCompiled(prog *CompiledProgram, input, expectedOutput string) *JudgeResult {
	if prog.IsJava {
		return j.runJavaCompiled(prog.JavaClassDir, input, expectedOutput, prog.TmpDir)
	}
	return j.runBinary(prog.BinaryPath, input, expectedOutput, prog.TmpDir)
}

// Run 统一入口，根据语言分发（兼容旧接口，内部会编译+执行）
func (j *Judge) Run(language, code, input, expectedOutput, driverCode string) *JudgeResult {
	prog, result := j.Compile(language, code, driverCode)
	if result != nil {
		return result
	}
	defer prog.Cleanup()
	return j.RunCompiled(prog, input, expectedOutput)
}

// ===== 编译方法 =====

func (j *Judge) compileGo(code, driverCode string) (*CompiledProgram, *JudgeResult) {
	tmpDir, err := os.MkdirTemp("", "oj_go_*")
	if err != nil {
		return nil, &JudgeResult{Status: "System Error", Error: fmt.Sprintf("Failed to create temp dir: %v", err)}
	}

	// 拼接最终代码
	finalCode := code
	if driverCode != "" {
		insertPos := findImportEnd(driverCode)
		if insertPos >= 0 {
			finalCode = driverCode[:insertPos] + "\n" + code + "\n" + driverCode[insertPos:]
		} else {
			pkgEnd := strings.Index(driverCode, "\n")
			if pkgEnd >= 0 {
				finalCode = driverCode[:pkgEnd+1] + "\n" + code + "\n" + driverCode[pkgEnd+1:]
			} else {
				finalCode = driverCode + "\n" + code
			}
		}
	}

	// 自动补充用户代码中使用但 import 块中缺失的标准库包
	finalCode = autoFixGoImports(finalCode)

	codeFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(codeFile, []byte(finalCode), 0644); err != nil {
		os.RemoveAll(tmpDir)
		return nil, &JudgeResult{Status: "System Error", Error: fmt.Sprintf("Failed to write code file: %v", err)}
	}

	binaryFile := filepath.Join(tmpDir, "main")
	compileCmd := exec.Command(j.goPath, "build", "-o", binaryFile, codeFile)
	compileCmd.Dir = tmpDir
	var compileErr bytes.Buffer
	compileCmd.Stderr = &compileErr
	if err := compileCmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return nil, &JudgeResult{Status: "Compile Error", Error: compileErr.String()}
	}

	return &CompiledProgram{TmpDir: tmpDir, BinaryPath: binaryFile}, nil
}

func (j *Judge) compileC(code string) (*CompiledProgram, *JudgeResult) {
	tmpDir, err := os.MkdirTemp("", "oj_c_*")
	if err != nil {
		return nil, &JudgeResult{Status: "System Error", Error: fmt.Sprintf("Failed to create temp dir: %v", err)}
	}

	codeFile := filepath.Join(tmpDir, "main.c")
	if err := os.WriteFile(codeFile, []byte(code), 0644); err != nil {
		os.RemoveAll(tmpDir)
		return nil, &JudgeResult{Status: "System Error", Error: fmt.Sprintf("Failed to write code file: %v", err)}
	}

	binaryFile := filepath.Join(tmpDir, "main")
	compileCmd := exec.Command(j.gccPath, "-O2", "-o", binaryFile, codeFile, "-lm")
	var compileErr bytes.Buffer
	compileCmd.Stderr = &compileErr
	if err := compileCmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return nil, &JudgeResult{Status: "Compile Error", Error: compileErr.String()}
	}

	return &CompiledProgram{TmpDir: tmpDir, BinaryPath: binaryFile}, nil
}

func (j *Judge) compileCpp(code string) (*CompiledProgram, *JudgeResult) {
	tmpDir, err := os.MkdirTemp("", "oj_cpp_*")
	if err != nil {
		return nil, &JudgeResult{Status: "System Error", Error: fmt.Sprintf("Failed to create temp dir: %v", err)}
	}

	codeFile := filepath.Join(tmpDir, "main.cpp")
	if err := os.WriteFile(codeFile, []byte(code), 0644); err != nil {
		os.RemoveAll(tmpDir)
		return nil, &JudgeResult{Status: "System Error", Error: fmt.Sprintf("Failed to write code file: %v", err)}
	}

	binaryFile := filepath.Join(tmpDir, "main")
	compileCmd := exec.Command(j.gppPath, "-O2", "-std=c++17", "-o", binaryFile, codeFile, "-lm")
	var compileErr bytes.Buffer
	compileCmd.Stderr = &compileErr
	if err := compileCmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return nil, &JudgeResult{Status: "Compile Error", Error: compileErr.String()}
	}

	return &CompiledProgram{TmpDir: tmpDir, BinaryPath: binaryFile}, nil
}

func (j *Judge) compileJava(code string) (*CompiledProgram, *JudgeResult) {
	tmpDir, err := os.MkdirTemp("", "oj_java_*")
	if err != nil {
		return nil, &JudgeResult{Status: "System Error", Error: fmt.Sprintf("Failed to create temp dir: %v", err)}
	}

	codeFile := filepath.Join(tmpDir, "Main.java")
	if err := os.WriteFile(codeFile, []byte(code), 0644); err != nil {
		os.RemoveAll(tmpDir)
		return nil, &JudgeResult{Status: "System Error", Error: fmt.Sprintf("Failed to write code file: %v", err)}
	}

	compileCmd := exec.Command(j.javacPath, codeFile)
	compileCmd.Dir = tmpDir
	var compileErr bytes.Buffer
	compileCmd.Stderr = &compileErr
	if err := compileCmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return nil, &JudgeResult{Status: "Compile Error", Error: compileErr.String()}
	}

	return &CompiledProgram{TmpDir: tmpDir, IsJava: true, JavaClassDir: tmpDir}, nil
}

// ===== 执行方法 =====

// runBinary 运行已编译的二进制文件（Go/C/C++）
// 沙盒隔离：通过 ulimit 限制子进程的文件大小、进程数、内存（作为双保险）
func (j *Judge) runBinary(binaryFile, input, expectedOutput, tmpDir string) *JudgeResult {
	timeOutputFile := filepath.Join(tmpDir, "time_output.txt")

	// 用 bash -c + ulimit 包裹运行命令，形成进程级沙盒：
	//   -f 16384  : 最大写文件 8 MB（512-byte blocks）防止磁盘炸弹
	//   -u 64     : 最大子进程数 64，防止 fork 炸弹
	//   -n 32     : 最大文件描述符，限制网络连接
	// 注意：不使用 ulimit -v 限制虚拟内存，因为 Go/Java 运行时启动时需要预留大量虚拟地址空间
	// 内存限制改为通过 /usr/bin/time -v 的 Maximum RSS 事后检查
	sandboxCmd := fmt.Sprintf(
		"ulimit -f 16384 -u 64 -n 32 2>/dev/null; exec /usr/bin/time -v -o %s %s",
		timeOutputFile, binaryFile,
	)
	runCmd := exec.Command("/bin/bash", "-c", sandboxCmd)
	runCmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	runCmd.Stdout = &stdout
	runCmd.Stderr = &stderr

	// 设置进程组，确保子进程树可以一起终止
	runCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// 墙钟时间仅用于超时检测
	wallStart := time.Now()
	err := runCmd.Start()
	if err != nil {
		return &JudgeResult{Status: "System Error", Error: fmt.Sprintf("Failed to start process: %v", err)}
	}

	// 使用 timer 在超时时整体 kill 进程组
	timer := time.AfterFunc(time.Duration(j.timeLimit)*time.Millisecond+500*time.Millisecond, func() {
		if runCmd.Process != nil {
			syscall.Kill(-runCmd.Process.Pid, syscall.SIGKILL)
		}
	})
	runErr := runCmd.Wait()
	timer.Stop()

	wallElapsed := time.Since(wallStart)
	wallTimeMs := int(wallElapsed.Milliseconds())

	// 从 time -v 解析 CPU 时间和内存
	usage := parseResourceUsage(timeOutputFile)
	timeUsed := usage.CPUTimeMs
	memoryUsed := usage.MemoryKB

	// 内存超限检查（基于实际 RSS）
	memoryLimitKB := j.memoryLimit * 1024
	if memoryUsed > memoryLimitKB {
		return &JudgeResult{Status: "Memory Limit Exceeded", TimeUsed: timeUsed, MemoryUsed: memoryUsed}
	}

	if runErr != nil {
		// 超时检测：墙钟时间超过限制视为超时
		if strings.Contains(runErr.Error(), "signal:") || wallTimeMs >= j.timeLimit {
			return &JudgeResult{Status: "Time Limit Exceeded", TimeUsed: j.timeLimit, MemoryUsed: memoryUsed}
		}
		return &JudgeResult{Status: "Runtime Error", Error: stderr.String(), TimeUsed: timeUsed, MemoryUsed: memoryUsed}
	}

	output := strings.TrimSpace(stdout.String())
	expected := strings.TrimSpace(expectedOutput)

	if expected == "" || output == expected {
		return &JudgeResult{Status: "Accepted", Output: output, TimeUsed: timeUsed, MemoryUsed: memoryUsed}
	}
	return &JudgeResult{Status: "Wrong Answer", Output: output, Expected: expected, TimeUsed: timeUsed, MemoryUsed: memoryUsed}
}

// runJavaCompiled 运行已编译的 Java 程序
// 沙盒隔离：同样通过 ulimit + Setpgid 实现进程级资源限制
func (j *Judge) runJavaCompiled(classDir, input, expectedOutput, tmpDir string) *JudgeResult {
	timeOutputFile := filepath.Join(tmpDir, "time_output.txt")

	sandboxCmd := fmt.Sprintf(
		"ulimit -f 16384 -u 64 -n 64 2>/dev/null; exec /usr/bin/time -v -o %s %s -cp %s -Xmx%dm Main",
		timeOutputFile, j.javaPath, classDir, j.memoryLimit,
	)
	runCmd := exec.Command("/bin/bash", "-c", sandboxCmd)
	runCmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	runCmd.Stdout = &stdout
	runCmd.Stderr = &stderr
	runCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// 墙钟时间仅用于超时检测
	wallStart := time.Now()
	err := runCmd.Start()
	if err != nil {
		return &JudgeResult{Status: "System Error", Error: fmt.Sprintf("Failed to start process: %v", err)}
	}

	timer := time.AfterFunc(time.Duration(j.timeLimit)*time.Millisecond+500*time.Millisecond, func() {
		if runCmd.Process != nil {
			syscall.Kill(-runCmd.Process.Pid, syscall.SIGKILL)
		}
	})
	runErr := runCmd.Wait()
	timer.Stop()

	wallElapsed := time.Since(wallStart)
	wallTimeMs := int(wallElapsed.Milliseconds())

	// 从 time -v 解析 CPU 时间和内存
	usage := parseResourceUsage(timeOutputFile)
	timeUsed := usage.CPUTimeMs
	memoryUsed := usage.MemoryKB

	// 内存超限检查（基于实际 RSS）
	memoryLimitKB := j.memoryLimit * 1024
	if memoryUsed > memoryLimitKB {
		return &JudgeResult{Status: "Memory Limit Exceeded", TimeUsed: timeUsed, MemoryUsed: memoryUsed}
	}

	if runErr != nil {
		// 超时检测：墙钟时间超过限制视为超时
		if strings.Contains(runErr.Error(), "signal:") || wallTimeMs >= j.timeLimit {
			return &JudgeResult{Status: "Time Limit Exceeded", TimeUsed: j.timeLimit, MemoryUsed: memoryUsed}
		}
		return &JudgeResult{Status: "Runtime Error", Error: stderr.String(), TimeUsed: timeUsed, MemoryUsed: memoryUsed}
	}

	output := strings.TrimSpace(stdout.String())
	expected := strings.TrimSpace(expectedOutput)
	if expected == "" || output == expected {
		return &JudgeResult{Status: "Accepted", Output: output, TimeUsed: timeUsed, MemoryUsed: memoryUsed}
	}
	return &JudgeResult{Status: "Wrong Answer", Output: output, Expected: expected, TimeUsed: timeUsed, MemoryUsed: memoryUsed}
}

// ===== 兼容旧接口的方法（单独编译+执行） =====

// RunGo 执行Go代码（兼容旧接口）
func (j *Judge) RunGo(code, input string, expectedOutput string, driverCode ...string) *JudgeResult {
	dc := ""
	if len(driverCode) > 0 {
		dc = driverCode[0]
	}
	return j.Run("Go", code, input, expectedOutput, dc)
}

// RunC 编译并执行C代码（兼容旧接口）
func (j *Judge) RunC(code, input, expectedOutput string) *JudgeResult {
	return j.Run("C", code, input, expectedOutput, "")
}

// RunCpp 编译并执行C++代码（兼容旧接口）
func (j *Judge) RunCpp(code, input, expectedOutput string) *JudgeResult {
	return j.Run("C++", code, input, expectedOutput, "")
}

// RunJava 编译并执行Java代码（兼容旧接口）
func (j *Judge) RunJava(code, input, expectedOutput string) *JudgeResult {
	return j.Run("Java", code, input, expectedOutput, "")
}

// ===== 辅助函数 =====

// ResourceUsage 资源使用统计
type ResourceUsage struct {
	CPUTimeMs   int // CPU 时间（毫秒）= User time + System time
	MemoryKB    int // 内存峰值（KB）
	WallTimeMs  int // 墙钟时间（毫秒），仅用于超时检测
}

// parseResourceUsage 从 /usr/bin/time -v 的输出中解析 CPU 时间和内存使用
// CPU 时间 = User time + System time，比墙钟时间更稳定
func parseResourceUsage(filename string) ResourceUsage {
	data, err := os.ReadFile(filename)
	if err != nil {
		return ResourceUsage{}
	}

	var userTime, sysTime float64
	var memoryKB int

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		
		// 解析 User time (seconds): 0.01
		if strings.Contains(line, "User time (seconds):") {
			parts := strings.Fields(line)
			if len(parts) >= 1 {
				if val, err := strconv.ParseFloat(parts[len(parts)-1], 64); err == nil {
					userTime = val
				}
			}
		}
		
		// 解析 System time (seconds): 0.00
		if strings.Contains(line, "System time (seconds):") {
			parts := strings.Fields(line)
			if len(parts) >= 1 {
				if val, err := strconv.ParseFloat(parts[len(parts)-1], 64); err == nil {
					sysTime = val
				}
			}
		}
		
		// 解析 Maximum resident set size (kbytes): 1234
		if strings.Contains(line, "Maximum resident set size") {
			parts := strings.Fields(line)
			if len(parts) >= 1 {
				if val, err := strconv.Atoi(parts[len(parts)-1]); err == nil {
					memoryKB = val
				}
			}
		}
	}

	// CPU 时间 = User + System，转换为毫秒
	cpuTimeMs := int((userTime + sysTime) * 1000)

	return ResourceUsage{
		CPUTimeMs:  cpuTimeMs,
		MemoryKB:   memoryKB,
	}
}

// autoFixGoImports 扫描代码中使用的标准库包，自动补充到 import 块中
func autoFixGoImports(code string) string {
	// 常见标准库包及其使用特征（包名.函数/类型）
	stdPkgs := map[string]*regexp.Regexp{
		"fmt":            regexp.MustCompile(`\bfmt\.`),
		"strconv":        regexp.MustCompile(`\bstrconv\.`),
		"strings":        regexp.MustCompile(`\bstrings\.`),
		"sort":           regexp.MustCompile(`\bsort\.`),
		"math":           regexp.MustCompile(`\bmath\.`),
		"math/big":       regexp.MustCompile(`\bbig\.`),
		"math/bits":      regexp.MustCompile(`\bbits\.`),
		"math/rand":      regexp.MustCompile(`\brand\.`),
		"bufio":          regexp.MustCompile(`\bbufio\.`),
		"bytes":          regexp.MustCompile(`\bbytes\.`),
		"os":             regexp.MustCompile(`\bos\.`),
		"io":             regexp.MustCompile(`\bio\.`),
		"errors":         regexp.MustCompile(`\berrors\.`),
		"regexp":         regexp.MustCompile(`\bregexp\.`),
		"sync":           regexp.MustCompile(`\bsync\.`),
		"unicode":        regexp.MustCompile(`\bunicode\.`),
		"unicode/utf8":   regexp.MustCompile(`\butf8\.`),
		"container/heap": regexp.MustCompile(`\bheap\.`),
		"container/list": regexp.MustCompile(`\blist\.`),
	}

	// 找出 import 块中已有的包
	existingImports := make(map[string]bool)
	importRe := regexp.MustCompile(`"([^"]+)"`)
	// 找 import 区域
	importStart := strings.Index(code, "import (")
	importEnd := -1
	if importStart >= 0 {
		importEnd = strings.Index(code[importStart:], "\n)")
		if importEnd >= 0 {
			importEnd += importStart
			importBlock := code[importStart : importEnd+2]
			for _, m := range importRe.FindAllStringSubmatch(importBlock, -1) {
				existingImports[m[1]] = true
			}
		}
	}

	if importEnd < 0 {
		return code // 没有 import 块，不处理
	}

	// 扫描代码体（import 块之后的部分）中使用了哪些包
	codeBody := code[importEnd+2:]
	var missing []string
	for pkg, re := range stdPkgs {
		if !existingImports[pkg] && re.MatchString(codeBody) {
			missing = append(missing, pkg)
		}
	}

	if len(missing) == 0 {
		return code
	}

	// 在 import 块的 ) 前插入缺失的包
	insertPoint := importEnd + 1 // "\n)" 的 ")" 位置
	var additions string
	for _, pkg := range missing {
		additions += fmt.Sprintf("\t\"%s\"\n", pkg)
	}

	return code[:insertPoint] + additions + code[insertPoint:]
}

// findImportEnd 找到 import 块结束位置
func findImportEnd(code string) int {
	closeIdx := strings.LastIndex(code, "\n)")
	if closeIdx >= 0 {
		return closeIdx + 2
	}
	lines := strings.Split(code, "\n")
	pos := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "import ") {
			pos += len(line) + 1
			return pos
		}
		pos += len(line) + 1
	}
	return -1
}
