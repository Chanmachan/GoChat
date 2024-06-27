function loginWithGoogle() {
    // Google OAuth リダイレクト処理
    window.location.href = 'https://accounts.google.com/o/oauth2/v2/auth?...';
}

function showLogin() {
    const app = document.getElementById('app');
    app.innerHTML = `
        <h1>Login to Chat Room</h1>
        <button onclick="loginWithGoogle()">Login with Google</button>
    `;
}
