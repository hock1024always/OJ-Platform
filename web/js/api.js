// API基础配置
// 自动检测当前主机地址
const API_BASE_URL = window.location.protocol + '//' + window.location.host + '/api/v1';

// API请求封装
class API {
    static async request(endpoint, options = {}) {
        const url = `${API_BASE_URL}${endpoint}`;
        const token = localStorage.getItem('token');

        const headers = {
            'Content-Type': 'application/json',
            ...options.headers,
        };

        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }

        try {
            const response = await fetch(url, {
                ...options,
                headers,
            });

            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.message || 'Request failed');
            }

            return data;
        } catch (error) {
            console.error('API Error:', error);
            throw error;
        }
    }

    // 用户相关
    static async register(username, email, password) {
        return this.request('/register', {
            method: 'POST',
            body: JSON.stringify({ username, email, password }),
        });
    }

    static async login(username, password) {
        return this.request('/login', {
            method: 'POST',
            body: JSON.stringify({ username, password }),
        });
    }

    static async getProfile() {
        return this.request('/profile');
    }

    // 题目相关
    static async getProblems(page = 1, pageSize = 20) {
        return this.request(`/problems?page=${page}&page_size=${pageSize}`);
    }

    static async getProblem(id) {
        return this.request(`/problems/${id}`);
    }

    static async createProblem(problemData) {
        return this.request('/problems', {
            method: 'POST',
            body: JSON.stringify(problemData),
        });
    }

    static async importProblem(problemData) {
        return this.request('/problems/import', {
            method: 'POST',
            body: JSON.stringify(problemData),
        });
    }

    // 判题相关
    static async submitCode(problemId, code, language) {
        return this.request('/submit', {
            method: 'POST',
            body: JSON.stringify({
                problem_id: problemId,
                code,
                language,
            }),
        });
    }

    static async getSubmission(id) {
        return this.request(`/submissions/${id}`);
    }

    // 运行测试（用户自定义输入）
    static async runTest(problemId, code, language, input) {
        return this.request('/test', {
            method: 'POST',
            body: JSON.stringify({
                problem_id: problemId,
                code,
                language,
                input,
            }),
        });
    }

    // 排行榜
    static async getGlobalLeaderboard() {
        return this.request('/leaderboard');
    }

    static async getProblemLeaderboard(problemId) {
        return this.request(`/problems/${problemId}/leaderboard`);
    }

    // 管理员
    static async getAdminSubmissions(page = 1, pageSize = 20, status = '', problemId = '') {
        let url = `/admin/submissions?page=${page}&pageSize=${pageSize}`;
        if (status) url += `&status=${encodeURIComponent(status)}`;
        if (problemId) url += `&problem_id=${problemId}`;
        return this.request(url);
    }

    static async getSubmissionCode(id) {
        return this.request(`/admin/submissions/${id}`);
    }
}

// 工具函数
function checkAuth() {
    const token = localStorage.getItem('token');
    if (!token) {
        window.location.href = '/index.html';
        return false;
    }
    return true;
}

function logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('username');
    window.location.href = '/index.html';
}

function showError(elementId, message) {
    const element = document.getElementById(elementId);
    if (element) {
        element.textContent = message;
        element.style.display = 'block';
    }
}

function clearError(elementId) {
    const element = document.getElementById(elementId);
    if (element) {
        element.textContent = '';
        element.style.display = 'none';
    }
}
