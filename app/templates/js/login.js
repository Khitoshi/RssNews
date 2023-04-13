// ログイン状態を管理する変数
let isLoggedIn = false;

// ログイン状態によって表示するリンクを変更する関数
function updateNavLinks() {
  /*
  const loginLink = document.getElementById('login-link');
  if (isLoggedIn) {
    loginLink.innerText = 'Logout';
    loginLink.removeEventListener('click', handleLoginClick);
    loginLink.addEventListener('click', handleLogoutClick);
  } else {
    loginLink.innerText = 'Login';
    loginLink.removeEventListener('click', handleLogoutClick);
    loginLink.addEventListener('click', handleLoginClick);
  }*/

  //
  if (document.cookie.indexOf('userId=') !== -1){
    document.getElementById("login-link").innerHTML = "<a href='/login' id='login-link'>Login</a>";
    //loginLink.innerText = 'Login';
    loginLink.removeEventListener('click', handleLogoutClick);
    loginLink.addEventListener('click', handleLoginClick);
  }else{//存在しなかった場合
    document.getElementById("logout-link").innerHTML = "<a href='/logout' id='logout-link'>Logout</a>";
    //loginLink.innerText = 'Logout';
    loginLink.removeEventListener('click', handleLoginClick);
    loginLink.addEventListener('click', handleLogoutClick);
  }
}

// ログインリンクがクリックされたときの処理
function handleLoginClick() {
  isLoggedIn = true;
  updateNavLinks();
  //TODO:このファイルでクッキーの確認をしてlogin | logoutの遷移を変更させる処理をさせる
}

// ログアウトリンクがクリックされたときの処理
function handleLogoutClick() {
  isLoggedIn = false;
  updateNavLinks();
}

// ページ読み込み時に初期化する処理
function initialize() {
  updateNavLinks();
}

initialize();