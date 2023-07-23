import { createContextMenu } from './contextMenu.js';

const groupContainer = document.getElementById("gridContainer");
let serviceRelations = {};
let currentState = "normal"

function alignMenuCorner(position, menuIcon, contextMenu) {
  const menuIconRect = menuIcon.getBoundingClientRect();
  if (position === "topRight") {
    const posY = menuIconRect.top;
    contextMenu.style.right = window.innerWidth - menuIconRect.right + "px";
    contextMenu.style.top = posY + "px";
  } else if (position === "topLeft") {
    const posX = menuIconRect.left;
    const posY = menuIconRect.top;
    contextMenu.style.left = posX + "px";
    contextMenu.style.top = posY + "px";
  } else if (position === "bottomRight") {
    contextMenu.style.right = window.innerWidth - menuIconRect.right + "px";
    contextMenu.style.bottom = window.innerHeight - menuIconRect.bottom + "px";
  } else if (position === "bottomLeft") {
    const posX = menuIconRect.left;
    contextMenu.style.left = posX + "px";
    contextMenu.style.bottom = window.innerHeight - menuIconRect.bottom + "px";
  }
}


function createGridItems(items, color) {
  const gridItems = [];
  items.forEach(item => {
    const gridItem = document.createElement("div");
    gridItem.className = "gridItem";
    // Set the background color from the data.yaml file
    gridItem.style.borderColor = color;
    gridItem.id = item.name; // Set the id attribute for each grid item (using the item name)

    const firstRow = document.createElement("div");
    firstRow.className = "firstRow";

    const titleElement = document.createElement("h3");
    titleElement.className = "title";
    titleElement.textContent = item.name;

    const menuIcon = document.createElement("i");
    menuIcon.className = "fas fa-bars menuIcon";
    
    // Append the title element and menu icon to the first row
    firstRow.appendChild(titleElement);
    firstRow.appendChild(menuIcon);

    const secondRow = document.createElement("div");
    secondRow.className = "secondRow";
    // Add the click event listener for the menu icon
    menuIcon.addEventListener("click", event => handleItemClick(event, item.menu, "topRight"));
    // Add the click event listener for the grid item (item itself)
    gridItem.addEventListener("click", event => handleGridItemClick(event, item));
    // Create clickable buttons for the second row
    if (item.buttons) {
      item.buttons.forEach(buttonData => {
        const button = document.createElement("button");
        button.className = "clickableButton";
        button.textContent = buttonData.name;
        button.style.backgroundColor = buttonData.color;
        // Add click event listener to the buttons
        button.addEventListener("click", event => handleItemClick(event, buttonData.items, "topRight"));
        // Append the button to the second row
        secondRow.appendChild(button);
      });
    }
    // Append the first row and second row to the grid item
    gridItem.appendChild(firstRow);
    gridItem.appendChild(secondRow);

    gridItems.push(gridItem);
  });

  return gridItems;
}

function handleGridItemClick(event, item) {
  const clickedItemName = item.name;
  if (currentState === "normal") {
    disableUnrelatedItems(clickedItemName, serviceRelations);
  } else {
    disableUnrelatedItems(null, serviceRelations);
  }
  event.stopPropagation();
  // You can implement other actions for the clicked item here
}

function disableUnrelatedItems(clickedItemName, serviceRelations) {
  const gridItems = document.querySelectorAll(".gridItem");
  gridItems.forEach(gridItem => {
    const itemName = gridItem.id;

    // If clickedItemName is null, it means we want to reset all items to their original state
    if (clickedItemName === null) {
      gridItem.classList.remove("disabled");
      currentState = "normal"
    } else {
      // Check if the item is related to the clicked item (either a backend or consumer)
      const relatedToClicked = serviceRelations[clickedItemName].backends.includes(itemName) ||
                               serviceRelations[clickedItemName].consumers.includes(itemName);
      // Disable unrelated items, except for the clicked item
      if (itemName !== clickedItemName) {
        gridItem.classList.toggle("disabled", !relatedToClicked);
      }
      currentState = "highlighted"
    }
  });
}

function handleItemClick(event, itemMenu, position) {
  event.stopPropagation(); // Prevent event bubbling

  const targetElement = event.target;
  const contextMenu = createContextMenu(itemMenu);

  // Use the passed arguments to align the menu corner
  alignMenuCorner(position, targetElement, contextMenu);

  setTimeout(() => {
    contextMenu.classList.add("show"); // Add the "show" class after a small delay
  }, 10);

  contextMenu.style.display = "block";

  // Remove the context menu when the cursor is away from it
  const removeContextMenu = (e) => {
    if (!contextMenu.contains(e.target)) {
      document.body.removeChild(contextMenu);
      document.removeEventListener("mousemove", removeContextMenu);
    }
  };
  document.addEventListener("mousemove", removeContextMenu);
}

function createServiceRelations(groupData) {
  const serviceRelations = {};
  groupData.forEach(group => {
    group.items.forEach(item => {
      serviceRelations[item.name] = {
        backends: item.backends ? item.backends : [],
        consumers: []
      };
    });
  });
  groupData.forEach(group => {
    group.items.forEach(item => {
      const backends = item.backends ? item.backends : [];
      backends.forEach(backend => {
        if (serviceRelations[backend]) {
          serviceRelations[backend].consumers.push(item.name);
        }
      });
    });
  });

  return serviceRelations;
}

fetch('data.yaml')
  .then(response => response.text())
  .then(yamlString => {
    const groupData = jsyaml.load(yamlString);
    serviceRelations = createServiceRelations(groupData);
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
  document.addEventListener('click', function(event) {
    if( currentState === "highlighted") {
      disableUnrelatedItems(null, serviceRelations);
    }
  });  
