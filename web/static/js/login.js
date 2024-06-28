function loginWithGoogle() {
    // Google OAuth リダイレクト処理
    window.location.href = "/auth/";
}

function showLogin() {
    const app = document.getElementById('app');
    app.innerHTML = `
        <h1>Login to Chat Room</h1>
        <button onclick="checkAuthAndRedirect()">Login with Google</button>
    `;
}

function getCookie(name) {
    let cookieArray = document.cookie.split(';');
    for(let i = 0; i < cookieArray.length; i++) {
        let cookie = cookieArray[i];
        while (cookie.charAt(0) === ' ') {
            cookie = cookie.substring(1);
        }
        if (cookie.indexOf(name + "=") === 0) {
            return cookie.substring(name.length + 1, cookie.length);
        }
    }
    return "";
}

function checkAuthAndRedirect() {
    const authCookie = getCookie("auth-session");
    if (authCookie) {
        console.log("Already login")
        navigate('/room'); // 認証済みの場合はルーム選択画面に遷移
    } else {
        loginWithGoogle();
    }
}
