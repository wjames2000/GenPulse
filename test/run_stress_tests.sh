#!/bin/bash

# GenPulse 压力测试脚本
# 用于测试系统在高负载下的性能和稳定性

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 显示帮助信息
show_help() {
    echo "GenPulse 压力测试脚本"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help          显示此帮助信息"
    echo "  -t, --type TYPE     测试类型 (memory, concurrent, pipeline, all)"
    echo "  -c, --count N       并发数量 (默认: 10)"
    echo "  -d, --duration N    测试持续时间(秒) (默认: 60)"
    echo "  -m, --memory N      内存测试大小(MB) (默认: 100)"
    echo "  -o, --output DIR    输出目录 (默认: ./stress_test_results)"
    echo "  -v, --verbose       显示详细输出"
    echo "  -p, --profile       生成性能分析文件"
    echo ""
    echo "示例:"
    echo "  $0 --type all       运行所有压力测试"
    echo "  $0 --type concurrent --count 50 运行50并发测试"
    echo "  $0 --type memory --memory 500   运行500MB内存测试"
}

# 清理函数
cleanup() {
    log_info "清理测试文件..."
    
    # 停止可能还在运行的测试进程
    pkill -f "stress_test" || true
    
    # 删除临时文件
    rm -rf /tmp/genpulse_stress_* || true
}

# 内存压力测试
run_memory_stress_test() {
    local memory_mb=$1
    local duration=$2
    local output_dir=$3
    
    log_info "开始内存压力测试 (${memory_mb}MB, ${duration}秒)..."
    
    # 创建测试目录
    TEST_DIR="$output_dir/memory"
    mkdir -p "$TEST_DIR"
    
    # 运行内存测试
    log_info "运行内存分配和释放测试..."
    
    # 编译并运行内存测试程序
    cat > /tmp/memory_stress_test.go << 'EOF'
package main

import (
    "fmt"
    "os"
    "runtime"
    "strconv"
    "time"
)

func main() {
    if len(os.Args) < 4 {
        fmt.Println("用法: memory_stress_test <内存MB> <持续时间秒> <输出文件>")
        os.Exit(1)
    }
    
    memoryMB, _ := strconv.Atoi(os.Args[1])
    durationSec, _ := strconv.Atoi(os.Args[2])
    outputFile := os.Args[3]
    
    // 打开输出文件
    f, err := os.Create(outputFile)
    if err != nil {
        fmt.Printf("创建输出文件失败: %v\n", err)
        os.Exit(1)
    }
    defer f.Close()
    
    // 记录开始信息
    startTime := time.Now()
    fmt.Fprintf(f, "内存压力测试开始\n")
    fmt.Fprintf(f, "目标内存: %d MB\n", memoryMB)
    fmt.Fprintf(f, "测试时长: %d 秒\n", durationSec)
    fmt.Fprintf(f, "开始时间: %s\n", startTime.Format("2006-01-02 15:04:05"))
    
    // 分配内存块
    blockSize := memoryMB * 1024 * 1024 / 10 // 分成10个块
    var blocks [][]byte
    
    fmt.Fprintf(f, "\n开始分配内存...\n")
    
    for i := 0; i < 10; i++ {
        block := make([]byte, blockSize)
        for j := 0; j < len(block); j += 4096 {
            block[j] = byte(j % 256) // 写入数据确保内存被实际分配
        }
        blocks = append(blocks, block)
        
        // 记录内存状态
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        fmt.Fprintf(f, "已分配 %d/%d MB, 系统内存: %.2f MB\n", 
            (i+1)*memoryMB/10, memoryMB,
            float64(m.Sys)/1024/1024)
    }
    
    fmt.Fprintf(f, "\n内存分配完成，保持 %d 秒...\n", durationSec)
    
    // 保持内存分配状态
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    endTime := startTime.Add(time.Duration(durationSec) * time.Second)
    
    for time.Now().Before(endTime) {
        select {
        case <-ticker.C:
            var m runtime.MemStats
            runtime.ReadMemStats(&m)
            fmt.Fprintf(f, "内存状态: 分配=%.2f MB, 系统=%.2f MB, Goroutines=%d\n",
                float64(m.Alloc)/1024/1024,
                float64(m.Sys)/1024/1024,
                runtime.NumGoroutine())
            
            // 执行一些内存操作
            for i := range blocks {
                for j := 0; j < len(blocks[i]); j += 8192 {
                    blocks[i][j] = blocks[i][j] + 1
                }
            }
        }
    }
    
    // 释放内存
    fmt.Fprintf(f, "\n释放内存...\n")
    blocks = nil
    runtime.GC()
    
    // 记录结束信息
    endTime = time.Now()
    duration := endTime.Sub(startTime)
    
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Fprintf(f, "\n测试完成\n")
    fmt.Fprintf(f, "结束时间: %s\n", endTime.Format("2006-01-02 15:04:05"))
    fmt.Fprintf(f, "总耗时: %v\n", duration)
    fmt.Fprintf(f, "最终内存状态: 分配=%.2f MB, 系统=%.2f MB\n",
        float64(m.Alloc)/1024/1024,
        float64(m.Sys)/1024/1024)
    fmt.Fprintf(f, "Goroutines: %d\n", runtime.NumGoroutine())
    
    fmt.Println("内存压力测试完成")
}
EOF
    
    go run /tmp/memory_stress_test.go "$memory_mb" "$duration" "$TEST_DIR/memory_test.log" 2>&1 | tee "$TEST_DIR/memory_test.console.log"
    
    # 检查测试结果
    if grep -q "内存压力测试完成" "$TEST_DIR/memory_test.console.log"; then
        log_success "内存压力测试完成"
        
        # 分析测试结果
        log_info "内存测试结果摘要:"
        tail -20 "$TEST_DIR/memory_test.log" | grep -E "(内存状态|测试完成|总耗时|最终内存)"
        
        return 0
    else
        log_error "内存压力测试失败"
        return 1
    fi
}

# 并发压力测试
run_concurrent_stress_test() {
    local concurrent_count=$1
    local duration=$2
    local output_dir=$3
    
    log_info "开始并发压力测试 (${concurrent_count}并发, ${duration}秒)..."
    
    # 创建测试目录
    TEST_DIR="$output_dir/concurrent"
    mkdir -p "$TEST_DIR"
    
    # 运行并发测试
    log_info "运行高并发任务测试..."
    
    # 编译并运行并发测试程序
    cat > /tmp/concurrent_stress_test.go << 'EOF'
package main

import (
    "fmt"
    "os"
    "strconv"
    "sync"
    "sync/atomic"
    "time"
)

func main() {
    if len(os.Args) < 4 {
        fmt.Println("用法: concurrent_stress_test <并发数> <持续时间秒> <输出文件>")
        os.Exit(1)
    }
    
    concurrentCount, _ := strconv.Atoi(os.Args[1])
    durationSec, _ := strconv.Atoi(os.Args[2])
    outputFile := os.Args[3]
    
    // 打开输出文件
    f, err := os.Create(outputFile)
    if err != nil {
        fmt.Printf("创建输出文件失败: %v\n", err)
        os.Exit(1)
    }
    defer f.Close()
    
    // 记录开始信息
    startTime := time.Now()
    fmt.Fprintf(f, "并发压力测试开始\n")
    fmt.Fprintf(f, "并发数量: %d\n", concurrentCount)
    fmt.Fprintf(f, "测试时长: %d 秒\n", durationSec)
    fmt.Fprintf(f, "开始时间: %s\n", startTime.Format("2006-01-02 15:04:05"))
    
    // 测试统计
    var totalTasks uint64
    var successfulTasks uint64
    var failedTasks uint64
    
    // 创建任务通道
    taskChan := make(chan int, concurrentCount*10)
    doneChan := make(chan bool, concurrentCount)
    
    // 启动worker
    var wg sync.WaitGroup
    
    for i := 0; i < concurrentCount; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            
            for taskID := range taskChan {
                taskStart := time.Now()
                
                // 模拟任务处理
                time.Sleep(time.Millisecond * time.Duration(10+(taskID%100)))
                
                // 模拟随机失败
                success := true
                if taskID%100 == 0 { // 1%的失败率
                    success = false
                    atomic.AddUint64(&failedTasks, 1)
                } else {
                    atomic.AddUint64(&successfulTasks, 1)
                }
                
                taskDuration := time.Since(taskStart)
                
                // 记录任务完成
                atomic.AddUint64(&totalTasks, 1)
                
                if taskID%1000 == 0 {
                    fmt.Fprintf(f, "Worker %d: 任务 %d %s, 耗时: %v\n", 
                        workerID, taskID, 
                        map[bool]string{true: "成功", false: "失败"}[success],
                        taskDuration)
                }
            }
            
            doneChan <- true
        }(i)
    }
    
    // 生成任务
    go func() {
        taskID := 0
        endTime := startTime.Add(time.Duration(durationSec) * time.Second)
        
        for time.Now().Before(endTime) {
            taskChan <- taskID
            taskID++
            
            // 控制任务生成速率
            if taskID%1000 == 0 {
                time.Sleep(time.Millisecond * 10)
            }
        }
        
        close(taskChan)
    }()
    
    // 监控进度
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    progressEndTime := startTime.Add(time.Duration(durationSec) * time.Second)
    
    for time.Now().Before(progressEndTime) {
        select {
        case <-ticker.C:
            currentTotal := atomic.LoadUint64(&totalTasks)
            currentSuccess := atomic.LoadUint64(&successfulTasks)
            currentFailed := atomic.LoadUint64(&failedTasks)
            
            elapsed := time.Since(startTime)
            tasksPerSec := float64(currentTotal) / elapsed.Seconds()
            
            fmt.Fprintf(f, "进度: 总任务=%d, 成功=%d, 失败=%d, 速率=%.1f 任务/秒\n",
                currentTotal, currentSuccess, currentFailed, tasksPerSec)
        }
    }
    
    // 等待所有worker完成
    wg.Wait()
    close(doneChan)
    
    // 记录结束信息
    endTime := time.Now()
    totalDuration := endTime.Sub(startTime)
    
    finalTotal := atomic.LoadUint64(&totalTasks)
    finalSuccess := atomic.LoadUint64(&successfulTasks)
    finalFailed := atomic.LoadUint64(&failedTasks)
    
    tasksPerSec := float64(finalTotal) / totalDuration.Seconds()
    successRate := float64(finalSuccess) / float64(finalTotal) * 100
    
    fmt.Fprintf(f, "\n测试完成\n")
    fmt.Fprintf(f, "结束时间: %s\n", endTime.Format("2006-01-02 15:04:05"))
    fmt.Fprintf(f, "总耗时: %v\n", totalDuration)
    fmt.Fprintf(f, "总任务数: %d\n", finalTotal)
    fmt.Fprintf(f, "成功任务: %d\n", finalSuccess)
    fmt.Fprintf(f, "失败任务: %d\n", finalFailed)
    fmt.Fprintf(f, "任务速率: %.1f 任务/秒\n", tasksPerSec)
    fmt.Fprintf(f, "成功率: %.2f%%\n", successRate)
    
    fmt.Printf("并发压力测试完成: %d 任务, %.1f 任务/秒, %.2f%% 成功率\n", 
        finalTotal, tasksPerSec, successRate)
}
EOF
    
    go run /tmp/concurrent_stress_test.go "$concurrent_count" "$duration" "$TEST_DIR/concurrent_test.log" 2>&1 | tee "$TEST_DIR/concurrent_test.console.log"
    
    # 检查测试结果
    if grep -q "并发压力测试完成" "$TEST_DIR/concurrent_test.console.log"; then
        log_success "并发压力测试完成"
        
        # 分析测试结果
        log_info "并发测试结果摘要:"
        tail -10 "$TEST_DIR/concurrent_test.log" | grep -E "(测试完成|总任务数|任务速率|成功率)"
        
        return 0
    else
        log_error "并发压力测试失败"
        return 1
    fi
}

# 流水线压力测试
run_pipeline_stress_test() {
    local concurrent_count=$1
    local duration=$2
    local output_dir=$3
    
    log_info "开始流水线压力测试 (${concurrent_count}并发, ${duration}秒)..."
    
    # 创建测试目录
    TEST_DIR="$output_dir/pipeline"
    mkdir -p "$TEST_DIR"
    
    # 运行流水线测试
    log_info "运行多流水线并发测试..."
    
    # 这里可以调用实际的流水线测试代码
    # 由于流水线测试需要完整的系统环境，这里先创建一个模拟测试
    
    cat > /tmp/pipeline_stress_test.go << 'EOF'
package main

import (
    "fmt"
    "os"
    "strconv"
    "sync"
    "sync/atomic"
    "time"
)

func main() {
    if len(os.Args) < 4 {
        fmt.Println("用法: pipeline_stress_test <并发数> <持续时间秒> <输出文件>")
        os.Exit(1)
    }
    
    pipelineCount, _ := strconv.Atoi(os.Args[1])
    durationSec, _ := strconv.Atoi(os.Args[2])
    outputFile := os.Args[3]
    
    // 打开输出文件
    f, err := os.Create(outputFile)
    if err != nil {
        fmt.Printf("创建输出文件失败: %v\n", err)
        os.Exit(1)
    }
    defer f.Close()
    
    // 记录开始信息
    startTime := time.Now()
    fmt.Fprintf(f, "流水线压力测试开始\n")
    fmt.Fprintf(f, "流水线数量: %d\n", pipelineCount)
    fmt.Fprintf(f, "测试时长: %d 秒\n", durationSec)
    fmt.Fprintf(f, "开始时间: %s\n", startTime.Format("2006-01-02 15:04:05"))
    
    // 测试统计
    var completedPipelines uint64
    var successfulPipelines uint64
    var failedPipelines uint64
    
    // 模拟流水线阶段
    pipelineStages := []string{
        "需求分析",
        "架构设计",
        "前端开发",
        "后端开发",
        "测试验证",
        "部署发布",
    }
    
    // 运行流水线
    var wg sync.WaitGroup
    results := make(chan string, pipelineCount)
    
    for i := 0; i < pipelineCount; i++ {
        wg.Add(1)
        go func(pipelineID int) {
            defer wg.Done()
            
            pipelineStart := time.Now()
            success := true
            
            fmt.Fprintf(f, "流水线 %d 开始\n", pipelineID)
            
            // 执行各个阶段
            for stageIndex, stage := range pipelineStages {
                stageStart := time.Now()
                
                // 模拟阶段执行时间
                stageDuration := time.Millisecond * time.Duration(100+(pipelineID%50)+(stageIndex*20))
                time.Sleep(stageDuration)
                
                // 模拟阶段失败
                if pipelineID%20 == 0 && stageIndex == 3 { // 5%的流水线在第4阶段失败
                    success = false
                    fmt.Fprintf(f, "流水线 %d: 阶段 '%s' 失败\n", pipelineID, stage)
                    break
                }
                
                stageElapsed := time.Since(stageStart)
                fmt.Fprintf(f, "流水线 %d: 阶段 '%s' 完成, 耗时: %v\n", 
                    pipelineID, stage, stageElapsed)
            }
            
            pipelineElapsed := time.Since(pipelineStart)
            
            if success {
                atomic.AddUint64(&successfulPipelines, 1)
                results <- fmt.Sprintf("流水线 %d 成功完成, 耗时: %v", pipelineID, pipelineElapsed)
            } else {
                atomic.AddUint64(&failedPipelines, 1)
                results <- fmt.Sprintf("流水线 %d 失败, 耗时: %v", pipelineID, pipelineElapsed)
            }
            
            atomic.AddUint64(&completedPipelines, 1)
        }(i)
    }
    
    // 监控进度
    ticker := time.NewTicker(2 * time.Second)
    defer ticker.Stop()
    
    endTime := startTime.Add(time.Duration(durationSec) * time.Second)
    
    go func() {
        for time.Now().Before(endTime) {
            select {
            case <-ticker.C:
                currentCompleted := atomic.LoadUint64(&completedPipelines)
                currentSuccess := atomic.LoadUint64(&successfulPipelines)
                currentFailed := atomic.LoadUint64(&failedPipelines)
                
                elapsed := time.Since(startTime)
                completionRate := float64(currentCompleted) / float64(pipelineCount) * 100
                
                fmt.Fprintf(f, "进度: 完成=%d/%d (%.1f%%), 成功=%d, 失败=%d\n",
                    currentCompleted, pipelineCount, completionRate,
                    currentSuccess, currentFailed)
            }
        }
    }()
    
    // 等待所有流水线完成或超时
    done := make(chan bool)
    go func() {
        wg.Wait()
        done <- true
    }()
    
    select {
    case <-done:
        // 所有流水线完成
    case <-time.After(time.Duration(durationSec) * time.Second):
        // 超时
        fmt.Fprintf(f, "测试超时，强制结束\n")
    }
    
    // 收集结果
    close(results)
    
    // 记录结束信息
    finalCompleted := atomic.LoadUint64(&completedPipelines)
    finalSuccess := atomic.LoadUint64(&successfulPipelines)
    finalFailed := atomic.LoadUint64(&failedPipelines)
    
    totalDuration := time.Since(startTime)
    completionRate := float64(finalCompleted) / float64(pipelineCount) * 100
    successRate := float64(finalSuccess) / float64(finalCompleted) * 100
    
    fmt.Fprintf(f, "\n测试完成\n")
    fmt.Fprintf(f, "结束时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
    fmt.Fprintf(f, "总耗时: %v\n", totalDuration)
    fmt.Fprintf(f, "完成流水线: %d/%d (%.1f%%)\n", finalCompleted, pipelineCount, completionRate)
    fmt.Fprintf(f, "成功流水线: %d\n", finalSuccess)
    fmt.Fprintf(f, "失败流水线: %d\n", finalFailed)
    fmt.Fprintf(f, "成功率: %.2f%%\n", successRate)
    
    // 输出结果摘要
    fmt.Fprintf(f, "\n结果摘要:\n")
    for result := range results {
        fmt.Fprintf(f, "%s\n", result)
    }
    
    fmt.Printf("流水线压力测试完成: %d/%d 完成, %.1f%% 成功率\n", 
        finalCompleted, pipelineCount, successRate)
}
EOF
    
    go run /tmp/pipeline_stress_test.go "$concurrent_count" "$duration" "$TEST_DIR/pipeline_test.log" 2>&1 | tee "$TEST_DIR/pipeline_test.console.log"
    
    # 检查测试结果
    if grep -q "流水线压力测试完成" "$TEST_DIR/pipeline_test.console.log"; then
        log_success "流水线压力测试完成"
        
        # 分析测试结果
        log_info "流水线测试结果摘要:"
        tail -15 "$TEST_DIR/pipeline_test.log" | grep -E "(测试完成|完成流水线|成功率)"
        
        return 0
    else
        log_error "流水线压力测试失败"
        return 1
    fi
}

# 生成测试报告
generate_test_report() {
    local output_dir=$1
    local start_time=$2
    local end_time=$3
    
    log_info "生成压力测试报告..."
    
    REPORT_FILE="$output_dir/stress_test_report.md"
    
    cat > "$REPORT_FILE" << EOF
# GenPulse 压力测试报告

## 测试概况
- 测试时间: $(date -r "$start_time" "+%Y-%m-%d %H:%M:%S") 到 $(date -r "$end_time" "+%Y-%m-%d %H:%M:%S")
- 总耗时: $((end_time - start_time)) 秒
- 测试环境: $(uname -s) $(uname -m)
- Go版本: $(go version)

## 测试结果摘要

EOF
    
    # 添加各个测试的结果
    for test_type in memory concurrent pipeline; do
        TEST_DIR="$output_dir/$test_type"
        
        if [ -d "$TEST_DIR" ]; then
            echo "### ${test_type} 测试" >> "$REPORT_FILE"
            echo "" >> "$REPORT_FILE"
            
            if [ -f "$TEST_DIR/${test_type}_test.log" ]; then
                # 提取关键结果
                case $test_type in
                    memory)
                        grep -E "(测试完成|总耗时|最终内存)" "$TEST_DIR/${test_type}_test.log" | tail -5 >> "$REPORT_FILE"
                        ;;
                    concurrent)
                        grep -E "(测试完成|总任务数|任务速率|成功率)" "$TEST_DIR/${test_type}_test.log" | tail -5 >> "$REPORT_FILE"
                        ;;
                    pipeline)
                        grep -E "(测试完成|完成流水线|成功率)" "$TEST_DIR/${test_type}_test.log" | tail -5 >> "$REPORT_FILE"
                        ;;
                esac
            else
                echo "测试未执行或失败" >> "$REPORT_FILE"
            fi
            
            echo "" >> "$REPORT_FILE"
        fi
    done
    
    # 添加系统资源使用情况
    echo "## 系统资源使用情况" >> "$REPORT_FILE"
    echo "" >> "$REPORT_FILE"
    echo '```' >> "$REPORT_FILE"
    echo "CPU信息:" >> "$REPORT_FILE"
    sysctl -n machdep.cpu.brand_string 2>/dev/null || echo "无法获取CPU信息" >> "$REPORT_FILE"
    echo "" >> "$REPORT_FILE"
    echo "内存信息:" >> "$REPORT_FILE"
    free -h 2>/dev/null || vm_stat 2>/dev/null || echo "无法获取内存信息" >> "$REPORT_FILE"
    echo '```' >> "$REPORT_FILE"
    
    log_success "压力测试报告已生成: $REPORT_FILE"
}

# 主函数
main() {
    # 默认参数
    TEST_TYPE="all"
    CONCURRENT_COUNT=10
    TEST_DURATION=60
    MEMORY_MB=100
    OUTPUT_DIR="./stress_test_results"
    VERBOSE=false
    PROFILE=false
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -t|--type)
                TEST_TYPE="$2"
                shift 2
                ;;
            -c|--count)
                CONCURRENT_COUNT="$2"
                shift 2
                ;;
            -d|--duration)
                TEST_DURATION="$2"
                shift 2
                ;;
            -m|--memory)
                MEMORY_MB="$2"
                shift 2
                ;;
            -o|--output)
                OUTPUT_DIR="$2"
                shift 2
                ;;
            -v|--verbose)
                VERBOSE=true
                shift
                ;;
            -p|--profile)
                PROFILE=true
                shift
                ;;
            *)
                log_error "未知选项: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # 验证参数
    if [ "$CONCURRENT_COUNT" -lt 1 ] || [ "$CONCURRENT_COUNT" -gt 1000 ]; then
        log_error "并发数量必须在1-1000之间"
        exit 1
    fi
    
    if [ "$TEST_DURATION" -lt 1 ] || [ "$TEST_DURATION" -gt 3600 ]; then
        log_error "测试持续时间必须在1-3600秒之间"
        exit 1
    fi
    
    if [ "$MEMORY_MB" -lt 1 ] || [ "$MEMORY_MB" -gt 4096 ]; then
        log_error "内存测试大小必须在1-4096MB之间"
        exit 1
    fi
    
    # 打印测试配置
    log_info "=== 压力测试配置 ==="
    log_info "测试类型: $TEST_TYPE"
    log_info "并发数量: $CONCURRENT_COUNT"
    log_info "测试时长: ${TEST_DURATION}秒"
    log_info "内存大小: ${MEMORY_MB}MB"
    log_info "输出目录: $OUTPUT_DIR"
    log_info "详细输出: $VERBOSE"
    log_info "性能分析: $PROFILE"
    log_info "=================="
    
    # 清理旧的测试文件
    cleanup
    
    # 创建输出目录
    mkdir -p "$OUTPUT_DIR"
    
    # 记录开始时间
    START_TIME=$(date +%s)
    
    # 执行测试
    FAILED_TESTS=0
    
    # 运行内存测试
    if [ "$TEST_TYPE" = "all" ] || [ "$TEST_TYPE" = "memory" ]; then
        if ! run_memory_stress_test "$MEMORY_MB" "$TEST_DURATION" "$OUTPUT_DIR"; then
            FAILED_TESTS=$((FAILED_TESTS + 1))
        fi
    fi
    
    # 运行并发测试
    if [ "$TEST_TYPE" = "all" ] || [ "$TEST_TYPE" = "concurrent" ]; then
        if ! run_concurrent_stress_test "$CONCURRENT_COUNT" "$TEST_DURATION" "$OUTPUT_DIR"; then
            FAILED_TESTS=$((FAILED_TESTS + 1))
        fi
    fi
    
    # 运行流水线测试
    if [ "$TEST_TYPE" = "all" ] || [ "$TEST_TYPE" = "pipeline" ]; then
        if ! run_pipeline_stress_test "$CONCURRENT_COUNT" "$TEST_DURATION" "$OUTPUT_DIR"; then
            FAILED_TESTS=$((FAILED_TESTS + 1))
        fi
    fi
    
    # 记录结束时间
    END_TIME=$(date +%s)
    
    # 生成测试报告
    generate_test_report "$OUTPUT_DIR" "$START_TIME" "$END_TIME"
    
    # 显示测试结果
    log_info "=== 压力测试结果 ==="
    log_info "总耗时: $((END_TIME - START_TIME)) 秒"
    
    if [ $FAILED_TESTS -eq 0 ]; then
        log_success "所有压力测试通过！"
    else
        log_error "$FAILED_TESTS 个测试失败"
        log_info "详细日志请查看: $OUTPUT_DIR/"
        exit 1
    fi
    
    exit 0
}

# 设置退出时清理
trap cleanup EXIT

# 运行主函数
main "$@"