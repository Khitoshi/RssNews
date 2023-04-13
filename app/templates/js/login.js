// ログイン状態を管理する変数
let isLoggedIn = false;

// ログイン状態によって表示するリンクを変更する関数
function updateNavLinks() {
  const loginLink = document.getElementById('login-link');
  if (isLoggedIn) {
    loginLink.innerText = 'Logout';
    loginLink.removeEventListener('click', handleLoginClick);
    loginLink.addEventListener('click', handleLogoutClick);
  } else {
    loginLink.innerText = 'Login';
    loginLink.removeEventListener('click', handleLogoutClick);
    loginLink.addEventListener('click', handleLoginClick);
  }
}

// ログインリンクがクリックされたときの処理
function handleLoginClick() {
  isLoggedIn = true;
  updateNavLinks();
  alert('ログインしました。');
}

// ログアウトリンクがクリックされたときの処理
function handleLogoutClick() {
  isLoggedIn = false;
  updateNavLinks();
  alert('ログアウトしました。');
}

// ページ読み込み時に初期化する処理
function initialize() {
  updateNavLinks();
}

initialize();