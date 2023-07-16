const groupContainer = document.getElementById("gridContainer");
function createGridItems(items) {
  const gridItems = [];
  items.forEach(item => {
    const gridItem = document.createElement("div");
    gridItem.className = "gridItem";

    // Create the content container
    const contentContainer = document.createElement("div");
    contentContainer.className = "contentContainer";

    // Add any other content or elements you want for the grid item
    const titleElement = document.createElement("h3");
    titleElement.textContent = item.title;
    contentContainer.appendChild(titleElement);

    const menuIcon = document.createElement("i");
    menuIcon.className = "fas fa-bars menuIcon";
    // Add event listener for menu icon click
    menuIcon.addEventListener("click", () => {
      toggleDropdownMenu(); // Call a function to toggle the dropdown menu
    });

    // Add the menu icon, content container, and dropdown menu to the grid item
    gridItem.appendChild(contentContainer);
    gridItem.appendChild(menuIcon);

    gridItems.push(gridItem);
  });
  return gridItems;
}

function toggleDropdownMenu() {
  console.log("context menu clicked")
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

      const gridItems = createGridItems(group.items); // Call the function to create grid items
      gridItems.forEach(gridItem => {
        gridContainerInner.appendChild(gridItem);
      });

      setGroup.appendChild(gridContainerInner);
      groupContainer.appendChild(setGroup);
    });
  })
  .catch(console.error);
