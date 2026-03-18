// 认证页面逻辑

// 显示标签页
function showTab(tabName) {
    // 切换按钮状态
    document.querySelectorAll('.tab-btn').forEach(btn => {
        btn.classList.remove('active');
    });
    event.target.classList.add('active');

    // 切换表单显示
    document.querySelectorAll('.auth-form').forEach(form => {
        form.classList.remove('active');
    });
    document.getElementById(`${tabName}-form`).classList.add('active');
}

// 登录处理
document.getElementById('login-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    clearError('login-error');

    const username = document.getElementById('login-username').value;
    const password = document.getElementById('login-password').value;

    try {
        const response = await API.login(username, password);
        localStorage.setItem('token', response.data.token);
        localStorage.setItem('username', username);
        window.location.href = '/problems.html';
    } catch (error) {
        showError('login-error', error.message);
    }
});

// 注册处理
document.getElementById('register-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    clearError('register-error');

    const username = document.getElementById('register-username').value;
    const email = document.getElementById('register-email').value;
    const password = document.getElementById('register-password').value;

    try {
        const response = await API.register(username, email, password);
        alert('注册成功！请登录');
        showTab('login');
        document.querySelector('.tab-btn').click(); // 激活登录标签
    } catch (error) {
        showError('register-error', error.message);
    }
});

// 检查是否已登录
if (localStorage.getItem('token')) {
    window.location.href = '/problems.html';
}
