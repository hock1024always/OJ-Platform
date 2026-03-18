// 题目列表页面逻辑

let currentPage = 1;
const pageSize = 20;

// 检查认证
if (!checkAuth()) {
    throw new Error('Not authenticated');
}

// 显示用户名
document.getElementById('username').textContent = localStorage.getItem('username');

// 加载题目列表
async function loadProblems() {
    try {
        const response = await API.getProblems(currentPage, pageSize);
        const problems = response.data.problems || [];

        const listElement = document.getElementById('problems-list');
        listElement.innerHTML = problems.map(problem => `
            <div class="problem-item" onclick="goToProblem(${problem.id})">
                <span class="problem-id">#${problem.id}</span>
                <span class="problem-title">${problem.title}</span>
                <span class="difficulty-badge difficulty-${problem.difficulty}">${problem.difficulty}</span>
            </div>
        `).join('');

        // 更新分页信息
        document.getElementById('page-info').textContent = `第 ${currentPage} 页`;
        document.getElementById('prev-btn').disabled = currentPage === 1;
        document.getElementById('next-btn').disabled = problems.length < pageSize;
    } catch (error) {
        alert('加载题目失败: ' + error.message);
    }
}

// 翻页
function loadPage(page) {
    if (page < 1) return;
    currentPage = page;
    loadProblems();
}

// 跳转到题目详情
function goToProblem(id) {
    window.location.href = `/problem.html?id=${id}`;
}

// 难度筛选
document.getElementById('difficulty-filter').addEventListener('change', (e) => {
    // 简单实现：重新加载并过滤
    loadProblems();
});

// 初始加载
loadProblems();
