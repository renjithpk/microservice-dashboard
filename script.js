
const groupContainer = document.getElementById("gridContainer");
function createContextMenu(menuItems) {
  var contextMenu = document.createElement('div');
  contextMenu.className = 'context-menu';

  var ul = document.createElement('ul');

  Object.entries(menuItems).forEach(([menuItem, link]) => {
    var li = document.createElement('li');
    var a = document.createElement('a');
    a.href = link;
    a.textContent = menuItem;
    li.appendChild(a);
    ul.appendChild(li);
  });

  contextMenu.appendChild(ul);
  document.body.appendChild(contextMenu);

  return contextMenu;
}

function createGridItems(items, color) {
  const gridItems = [];
  items.forEach(item => {
    const gridItem = document.createElement("div");
    gridItem.className = "gridItem";
    // Set the background color from the data.yaml file
    gridItem.style.borderColor = color;

    const contentContainer = document.createElement("div");
    contentContainer.className = "contentContainer";

    const titleElement = document.createElement("h3");
    titleElement.textContent = item.name;
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
      title.textContent = group.group;
      setGroup.appendChild(title);

      const gridContainerInner = document.createElement("div");
      gridContainerInner.className = "gridContainer";

      const gridItems = createGridItems(group.items, group.color);
      gridItems.forEach(gridItem => {
        gridContainerInner.appendChild(gridItem);
      });

      setGroup.appendChild(gridContainerInner);
      groupContainer.appendChild(setGroup);
    });
  })
  .catch(console.error);

  