function createContextMenu() {
  var contextMenu = document.createElement('div');
  contextMenu.className = 'context-menu';

  var ul = document.createElement('ul');

  var li1 = document.createElement('li');
  var a1 = document.createElement('a');
  a1.href = '#';
  a1.textContent = 'Menu Item 1';
  li1.appendChild(a1);
  ul.appendChild(li1);

  var li2 = document.createElement('li');
  var a2 = document.createElement('a');
  a2.href = '#';
  a2.textContent = 'Menu Item 2';
  li2.appendChild(a2);
  ul.appendChild(li2);

  var li3 = document.createElement('li');
  var a3 = document.createElement('a');
  a3.href = '#';
  a3.textContent = 'Menu Item 3';
  li3.appendChild(a3);
  ul.appendChild(li3);

  contextMenu.appendChild(ul);
  document.body.appendChild(contextMenu);

  return contextMenu;
}

const groupContainer = document.getElementById("gridContainer");

function createContextMenu(menuItems) {
  var contextMenu = document.createElement('div');
  contextMenu.className = 'context-menu';

  var ul = document.createElement('ul');

  menuItems.forEach(item => {
    var li = document.createElement('li');
    var a = document.createElement('a');
    a.href = '#';
    a.textContent = item;
    li.appendChild(a);
    ul.appendChild(li);
  });

  contextMenu.appendChild(ul);
  document.body.appendChild(contextMenu);

  return contextMenu;
}


function createGridItems(items, bgColor) {
  const gridItems = [];
  items.forEach(item => {
    const gridItem = document.createElement("div");
    gridItem.className = "gridItem";
    // Set the background color from the data.yaml file
    gridItem.style.borderColor = bgColor ;

    const contentContainer = document.createElement("div");
    contentContainer.className = "contentContainer";

    const titleElement = document.createElement("h3");
    titleElement.textContent = item.title;
    contentContainer.appendChild(titleElement);

    const menuIcon = document.createElement("i");
    menuIcon.className = "fas fa-bars menuIcon";

    menuIcon.addEventListener("click", (event) => {
      event.stopPropagation(); // Prevent event bubbling
      const contextMenu = createContextMenu(item.menu);
      setTimeout(() => {
        contextMenu.classList.add("show"); // Add the "show" class after a small delay
      }, 10);
      const menuIconRect = menuIcon.getBoundingClientRect();
      const posX = menuIconRect.right - contextMenu.offsetWidth;
      const posY = menuIconRect.top;

      contextMenu.style.right = window.innerWidth - menuIconRect.right + "px";
      contextMenu.style.top = posY + "px";
      contextMenu.style.display = "block";

      // Remove the context menu when the cursor is away from it
      const removeContextMenu = (e) => {
        if (!contextMenu.contains(e.target)) {
          document.body.removeChild(contextMenu);
          document.removeEventListener("mousemove", removeContextMenu);
        }
      };
      document.addEventListener("mousemove", removeContextMenu);
    });

    // Add the menu icon, content container, and dropdown menu to the grid item
    gridItem.appendChild(contentContainer);
    gridItem.appendChild(menuIcon);

    gridItems.push(gridItem);
  });
  return gridItems;
}

fetch('data.yaml')
  .then(response => response.text())
  .then(yamlString => {
    const groupData = jsyaml.load(yamlString);

    groupData.forEach(group => {
      const setGroup = document.createElement("div");
      setGroup.className = "setGroup";

      const title = document.createElement("h3");
      title.textContent = group.name;
      setGroup.appendChild(title);

      const gridContainerInner = document.createElement("div");
      gridContainerInner.className = "gridContainer";

      const gridItems = createGridItems(group.items, group.bgColor);
      gridItems.forEach(gridItem => {
        gridContainerInner.appendChild(gridItem);
      });

      setGroup.appendChild(gridContainerInner);
      groupContainer.appendChild(setGroup);
    });
  })
  .catch(console.error);
