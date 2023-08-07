import { createContextMenu } from './contextMenu.js';
import { getName } from './jsonQuery.js';

const groupContainer = document.getElementById("gridContainer");
let relations = {};
let currentState = "normal";

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

        // Check if both image and text are provided
        if (buttonData.image && buttonData.name) {
          // Create a container for the image and text (flexbox layout)
          const contentContainer = document.createElement("div");
          contentContainer.style.display = "flex";
          contentContainer.style.alignItems = "center";

          // Create the image element
          const image = document.createElement("img");
          image.src = buttonData.image;
          // Set any additional image styles here if needed

          // Create the text element
          const text = document.createElement("span");
          getName(buttonData.name)
            .then(result => text.textContent = result)
            .catch(error => console.error(error.message));
          // Set any additional text styles here if needed

          // Append the image and text elements to the content container
          contentContainer.appendChild(image);
          contentContainer.appendChild(text);

          // Append the content container to the button
          button.appendChild(contentContainer);
        } else if (buttonData.image) {
          // If only the image is provided, add the image element directly to the button
          const image = document.createElement("img");
          image.src = buttonData.image;
          // Set any additional image styles here if needed

          button.appendChild(image);
        } else if (buttonData.name){
          getName(buttonData.name)
            .then(result => button.textContent = result)
            .catch(error => console.error(error.message));
        }
        button.style.backgroundColor = buttonData.color;
        if(buttonData.items) {
          // Add click event listener to the buttons
          button.addEventListener("click", event => handleItemClick(event, buttonData.items, "topRight"));
        } else if (buttonData.link) {
          button.addEventListener("click", function(event){
            event.stopPropagation();
            window.open(buttonData.link, "_blank");
          })
        }
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
    disableUnrelatedItems(clickedItemName);
  } else {
    disableUnrelatedItems(null);
  }
  event.stopPropagation();
}

function disableUnrelatedItems(clickedItemName) {
  const gridItems = document.querySelectorAll(".gridItem");
  gridItems.forEach(gridItem => {
    const itemName = gridItem.id;
    if (clickedItemName === null) {
      gridItem.classList.remove("disabled");
      currentState = "normal";
    } else {
      const relationInfo = relations[clickedItemName];
      const isRelatedToClicked = relationInfo && (
        (relationInfo.backends && relationInfo.backends.includes(itemName)) ||
        (relationInfo.consumers && relationInfo.consumers.includes(itemName))
      );
      // Disable unrelated items, except for the clicked item
      if (itemName !== clickedItemName) {
        gridItem.classList.toggle("disabled", !isRelatedToClicked);
      }
      currentState = "highlighted";
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
    groupData.forEach(group => {
      group.items.forEach(item => {
        relations[item.name] = {
          backends: item.backends ? item.backends : [],
          consumers: item.consumers ? item.consumers : []
        };
      });
    });

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
  if (currentState === "highlighted") {
    disableUnrelatedItems(null);
  }
});

