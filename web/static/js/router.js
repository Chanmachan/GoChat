function router() {
    const route = window.location.pathname;
    switch(route) {
        case '/login':
            showLogin();
            break;
        case '/room':
            showRoomSelection();
            break;
        case '/chat':
            showChat();
            break;
        default:
            // showLogin();  // デフォルトのビュー
            showRoomSelection();
    }
}

function navigate(path) {
    history.pushState({}, '', path);
    router();
}

window.addEventListener('popstate', router); // ブラウザの戻る/進むをサポート
