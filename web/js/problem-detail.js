// 题目详情页面逻辑

let currentProblemId = null;
let editor = null;
let currentProblem = null;

// 各语言默认模板（当题目没有对应语言模板时使用）
const DEFAULT_TEMPLATES = {
    'Go': `func solution() {
    // 请在此实现你的代码
}`,
    'C': `#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int main() {
    // 请在此实现你的代码
    return 0;
}`,
    'C++': `#include <bits/stdc++.h>
using namespace std;

int main() {
    // 请在此实现你的代码
    return 0;
}`,
    'Java': `import java.util.*;
import java.io.*;

public class Main {
    public static void main(String[] args) {
        Scanner sc = new Scanner(System.in);
        // 请在此实现你的代码
    }
}`
};

// CodeMirror 语言模式映射
const CM_MODES = {
    'Go': 'go',
    'C': 'text/x-csrc',
    'C++': 'text/x-c++src',
    'Java': 'text/x-java'
};

// 检查认证
if (!checkAuth()) {
    throw new Error('Not authenticated');
}

// 显示用户名
document.getElementById('username').textContent = localStorage.getItem('username');

// 获取URL参数
function getUrlParam(name) {
    const urlParams = new URLSearchParams(window.location.search);
    return urlParams.get(name);
}

// 加载题目详情
async function loadProblem() {
    currentProblemId = getUrlParam('id');
    if (!currentProblemId) {
        alert('题目ID不存在');
        window.location.href = '/problems.html';
        return;
    }

    try {
        const response = await API.getProblem(currentProblemId);
        const problem = response.data.problem;
        const testCases = response.data.testCases || [];

        // 显示题目信息
        document.getElementById('problem-title').textContent = `${problem.id}. ${problem.title}`;
        document.getElementById('problem-difficulty').textContent = problem.difficulty;
        document.getElementById('problem-difficulty').className = `difficulty-badge difficulty-${problem.difficulty}`;
        document.getElementById('time-limit').textContent = problem.time_limit;
        document.getElementById('memory-limit').textContent = problem.memory_limit;
        document.getElementById('problem-description').innerHTML = problem.description.replace(/\n/g, '<br>');

        // 显示测试用例
        const testCasesHtml = testCases.map((tc, index) => `
            <div class="test-case">
                <div class="test-case-label">示例 ${index + 1}:</div>
                <div class="test-case-content">
                    <div><strong>输入:</strong> ${tc.input.replace(/\n/g, '<br>')}</div>
                    <div><strong>输出:</strong> ${tc.output}</div>
                </div>
            </div>
        `).join('');
        document.getElementById('test-cases-list').innerHTML = testCasesHtml;

        // 设置输入提示
        if (testCases.length > 0) {
            document.getElementById('input-hint').textContent = 
                '输入格式示例: ' + testCases[0].input.replace(/\n/g, ' → ').substring(0, 50) + '...';
        }

        // 保存题目信息
        currentProblem = problem;

        // 初始化代码编辑器，使用题目的函数模板（Go语言默认）
        initEditor('Go', problem.function_template);
    } catch (error) {
        alert('加载题目失败: ' + error.message);
    }
}

// 初始化代码编辑器
function initEditor(language, functionTemplate) {
    const textarea = document.getElementById('code-input');
    const mode = CM_MODES[language] || 'go';

    if (editor) {
        editor.setOption('mode', mode);
        editor.setValue(functionTemplate || DEFAULT_TEMPLATES[language] || '');
        return;
    }

    editor = CodeMirror.fromTextArea(textarea, {
        mode: mode,
        theme: 'default',
        lineNumbers: true,
        indentUnit: 4,
        tabSize: 4,
        indentWithTabs: true,
        lineWrapping: true,
    });

    editor.setValue(functionTemplate || DEFAULT_TEMPLATES[language] || '');
}

// 切换语言时更新编辑器模式和模板
function onLanguageChange() {
    const lang = document.getElementById('language-select').value;
    const mode = CM_MODES[lang] || 'go';
    editor.setOption('mode', mode);

    // Go语言使用题目提供的函数模板，其他语言用默认通用模板
    const template = (lang === 'Go' && currentProblem && currentProblem.function_template)
        ? currentProblem.function_template
        : DEFAULT_TEMPLATES[lang] || '';
    editor.setValue(template);
}

// 提交代码
async function submitCode() {
    const code = editor.getValue();
    const language = document.getElementById('language-select').value;    if (!code.trim()) {
        alert('请输入代码');
        return;
    }

    try {
        const response = await API.submitCode(parseInt(currentProblemId), code, language);
        const submissionId = response.data.id;

        // 显示模态框
        showResultModal(submissionId);
    } catch (error) {
        alert('提交失败: ' + error.message);
    }
}

// 显示结果模态框
function showResultModal(submissionId) {
    const modal = document.getElementById('result-modal');
    const content = document.getElementById('result-content');

    content.innerHTML = '<div class="result-pending">正在判题中...</div>';
    modal.style.display = 'block';

    // 轮询获取结果
    const interval = setInterval(async () => {
        try {
            const response = await API.getSubmission(submissionId);
            const submission = response.data;

            if (submission.status !== 'Pending') {
                clearInterval(interval);
                displayResult(submission);
            }
        } catch (error) {
            clearInterval(interval);
            content.innerHTML = `<div class="result-wrong">查询失败: ${error.message}</div>`;
        }
    }, 1000);
}

// 显示判题结果
function displayResult(submission) {
    const content = document.getElementById('result-content');

    let statusClass = '';
    let statusText = '';

    switch (submission.status) {
        case 'Accepted':
            statusClass = 'result-accepted';
            statusText = '✓ 通过';
            break;
        case 'Wrong Answer':
            statusClass = 'result-wrong';
            statusText = '✗ 答案错误';
            break;
        case 'Compile Error':
            statusClass = 'result-wrong';
            statusText = '✗ 编译错误';
            break;
        case 'Time Limit Exceeded':
            statusClass = 'result-wrong';
            statusText = '✗ 超时';
            break;
        default:
            statusClass = 'result-wrong';
            statusText = submission.status;
    }

    // 性能可视化
    const timeUsed = submission.time_used || 0;
    const memoryUsed = submission.memory_used || 0;
    const timeLimit = 5000; // ms
    const memoryLimit = 256 * 1024; // KB (256MB)
    
    const timePercent = Math.min(100, (timeUsed / timeLimit) * 100);
    const memoryPercent = Math.min(100, (memoryUsed / memoryLimit) * 100);
    
    const memoryDisplay = memoryUsed >= 1024 
        ? (memoryUsed / 1024).toFixed(2) + ' MB'
        : memoryUsed + ' KB';

    let html = `
        <div class="${statusClass}" style="font-size: 1.5rem; margin-bottom: 1rem;">
            ${statusText}
        </div>
        <div class="result-details">
            <div class="perf-item" style="margin-bottom: 1rem;">
                <div style="display: flex; justify-content: space-between; margin-bottom: 0.3rem;">
                    <strong>执行时间</strong>
                    <span>${timeUsed} ms</span>
                </div>
                <div style="background: #e9ecef; border-radius: 4px; height: 8px; overflow: hidden;">
                    <div style="background: ${timePercent > 80 ? '#dc3545' : timePercent > 50 ? '#ffc107' : '#28a745'}; 
                                height: 100%; width: ${timePercent}%; transition: width 0.3s;"></div>
                </div>
                <div style="font-size: 12px; color: #666; margin-top: 2px;">限制: ${timeLimit} ms</div>
            </div>
            <div class="perf-item" style="margin-bottom: 1rem;">
                <div style="display: flex; justify-content: space-between; margin-bottom: 0.3rem;">
                    <strong>内存消耗</strong>
                    <span>${memoryDisplay}</span>
                </div>
                <div style="background: #e9ecef; border-radius: 4px; height: 8px; overflow: hidden;">
                    <div style="background: ${memoryPercent > 80 ? '#dc3545' : memoryPercent > 50 ? '#ffc107' : '#17a2b8'}; 
                                height: 100%; width: ${memoryPercent}%; transition: width 0.3s;"></div>
                </div>
                <div style="font-size: 12px; color: #666; margin-top: 2px;">限制: 256 MB</div>
            </div>
    `;

    if (submission.result) {
        html += `<p><strong>详细信息:</strong></p><pre style="background:#f8f9fa;padding:0.5rem;border-radius:4px;font-size:12px;max-height:150px;overflow-y:auto;">${submission.result}</pre>`;
    }

    html += '</div>';
    content.innerHTML = html;
}

// 关闭模态框
function closeModal() {
    document.getElementById('result-modal').style.display = 'none';
}

// 切换测试控制台标签
function switchConsoleTab(tab) {
    document.getElementById('tab-input').classList.toggle('active', tab === 'input');
    document.getElementById('tab-output').classList.toggle('active', tab === 'output');
    document.getElementById('pane-input').classList.toggle('active', tab === 'input');
    document.getElementById('pane-output').classList.toggle('active', tab === 'output');
}

// 运行自定义测试
async function runCustomTest() {
    const code = editor.getValue();
    const language = document.getElementById('language-select').value;
    const input = document.getElementById('test-input').value;

    if (!code.trim()) {
        alert('请输入代码');
        return;
    }

    // 切换到输出 pane 并显示"运行中"
    switchConsoleTab('output');
    const outputBox = document.getElementById('console-output-box');
    outputBox.className = 'console-output running';
    outputBox.textContent = '正在运行测试...';
    document.getElementById('console-perf').style.display = 'none';

    try {
        const response = await API.runTest(parseInt(currentProblemId), code, language, input);
        displayTestResult(response.data);
    } catch (error) {
        outputBox.className = 'console-output error';
        outputBox.textContent = '测试失败: ' + error.message;
    }
}

// 显示测试结果（内联控制台版本）
function displayTestResult(result) {
    const outputBox = document.getElementById('console-output-box');
    const perfSection = document.getElementById('console-perf');

    const isSuccess = result.status === 'Accepted' || result.status === 'Finished';
    const isCompileError = result.status === 'Compile Error';

    outputBox.className = 'console-output ' + (isSuccess ? 'success' : 'error');

    let text = '';
    if (isSuccess) {
        text = result.output || '(无输出)';
    } else if (isCompileError) {
        text = '[编译错误]\n' + (result.error || '');
    } else if (result.status === 'Time Limit Exceeded') {
        text = '[超时]';
    } else if (result.status === 'Runtime Error') {
        text = '[运行错误]\n' + (result.error || '');
    } else {
        text = '[' + result.status + ']\n' + (result.error || result.output || '');
    }
    outputBox.textContent = text;

    // 性能条
    if (isSuccess && (result.time_used || result.memory_used)) {
        const timeUsed = result.time_used || 0;
        const memoryUsed = result.memory_used || 0;
        const timeLimit = 5000;
        const memoryLimit = 256 * 1024;

        const timePercent = Math.min(100, (timeUsed / timeLimit) * 100);
        const memoryPercent = Math.min(100, (memoryUsed / memoryLimit) * 100);

        const memoryDisplay = memoryUsed >= 1024
            ? (memoryUsed / 1024).toFixed(2) + ' MB'
            : memoryUsed + ' KB';

        document.getElementById('perf-time-fill').style.width = timePercent + '%';
        document.getElementById('perf-time-fill').style.background = timePercent > 80 ? '#dc3545' : timePercent > 50 ? '#ffc107' : '#28a745';
        document.getElementById('perf-time-val').textContent = timeUsed + ' ms';

        document.getElementById('perf-mem-fill').style.width = memoryPercent + '%';
        document.getElementById('perf-mem-fill').style.background = memoryPercent > 80 ? '#dc3545' : memoryPercent > 50 ? '#ffc107' : '#17a2b8';
        document.getElementById('perf-mem-val').textContent = memoryDisplay;

        perfSection.style.display = 'block';
    } else {
        perfSection.style.display = 'none';
    }
}



// 点击模态框外部关闭
window.onclick = function(event) {
    const modal = document.getElementById('result-modal');
    if (event.target === modal) {
        closeModal();
    }
};

// 初始加载
loadProblem();
