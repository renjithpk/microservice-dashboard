import { getName } from './jsonQuery.js';

function createContextMenu(menuItems) {
  var contextMenu = document.createElement('div');
  contextMenu.className = 'context-menu';

  var ul = document.createElement('ul');

  menuItems.forEach((menuItem) => {
    var li = document.createElement('li');
    if(menuItem.link) {
      var a = document.createElement('a');
      a.href = menuItem.link;
      getName(menuItem.name)
        .then(result => a.textContent = result)
        .catch(error => console.error(error.message));
      a.target = '_blank'; // Open link in a new tab
      li.appendChild(a);
    } else {
      getName(menuItem.name)
        .then(result => li.textContent = result)
        .catch(error => console.error(error.message));
    }
    ul.appendChild(li);
  });

  contextMenu.appendChild(ul);
  document.body.appendChild(contextMenu);

  return contextMenu;
}

  
// Export the createContextMenu function
export { createContextMenu };
  