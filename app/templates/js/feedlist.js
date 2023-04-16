//const btn = document.querySelector('button');
const form = document.getElementById('form-block');
const closeIcons = document.querySelectorAll('.close-icon');
const items = document.querySelectorAll('.item');

// 新しいフォームを作成し、作成したフォームを返す関数
/*
function createNewForm() {
  //const newDiv = document.createElement('div');
  const newDiv = document.createElement('.item');
  newDiv.classList.add('item');

  const newForm = document.createElement('input');
  newForm.type = 'text';
  
  const newLabel = document.createElement('label');
  newLabel.textContent = '購読したいRSS-URL:';

  const newSpan = document.createElement('span');
  newSpan.classList.add('close-icon');
  newSpan.textContent = '✖';
  
  newLabel.appendChild(newForm);
  newDiv.appendChild(newLabel);
  newDiv.appendChild(newSpan);

  // 「✖」をクリックしたときの処理を追加
  newSpan.addEventListener('click', () => {
    newDiv.remove();
  });

  return newDiv;
}
*/

function createNewForm() {
  const newDiv = document.createElement('div');
  newDiv.classList.add('item');

  const newForm = document.createElement('input');
  newForm.type = 'text';

  const newLabel = document.createElement('label');
  newLabel.textContent = '購読したいRSS-URL:';

  const newSpan = document.createElement('span');
  newSpan.classList.add('close-icon');
  newSpan.textContent = '✖';

  newLabel.appendChild(newForm);
  newDiv.appendChild(newLabel);
  newDiv.appendChild(newSpan);

  // 「✖」をクリックしたときの処理を追加
  newSpan.addEventListener('click', () => {
    newDiv.remove();
  });

  const addButton = document.querySelector('.additem');
  addButton.parentNode.insertBefore(newDiv, addButton.nextSibling);

  return newDiv;
}



// 「✖」をクリックしたときの処理
for (let j = 0; j < closeIcons.length; j++) {
  closeIcons[j].addEventListener('click', () => {
    items[j].remove();
  });
}

//ボタンをクリックした時の動作
$('.additem').on('click', function() {//タイトル要素をクリックしたら
    form.appendChild(createNewForm());
});