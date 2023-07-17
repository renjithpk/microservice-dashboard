function createContextMenu(menuItems) {
  var contextMenu = document.createElement('div');
  contextMenu.className = 'context-menu';

  var ul = document.createElement('ul');

  menuItems.forEach((menuItem) => {
    var li = document.createElement('li');
    var a = document.createElement('a');
    a.href = menuItem.link;
    a.textContent = menuItem.name;
    a.target = '_blank'; // Open link in a new tab

    li.addEventListener('click', function () {
      a.click(); // Trigger the link click when the list item is clicked
    });

    li.appendChild(a);
    ul.appendChild(li);
  });

  contextMenu.appendChild(ul);
  document.body.appendChild(contextMenu);

  return contextMenu;
}

  
// Export the createContextMenu function
export { createContextMenu };
  