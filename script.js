const groupContainer = document.getElementById("gridContainer");

fetch('data.yaml')
  .then(response => response.text())
  .then(yamlString => {
    const data = jsyaml.load(yamlString);

    if (data?.data?.length > 0) {
      data.data.forEach(group => {
        const groupContainerDiv = document.createElement("div");
        groupContainerDiv.classList.add("grid-group");

        group.items.forEach(item => {
          const tile = document.createElement("div");
          tile.classList.add("grid-item");

          const menuIcon = document.createElement("i");
          menuIcon.classList.add("fas", "fa-bars");
          tile.appendChild(menuIcon);

          const titleElement = document.createElement("h3");
          titleElement.textContent = item.title;
          tile.appendChild(titleElement);

          const contentElement = document.createElement("p");
          contentElement.textContent = item.content;
          tile.appendChild(contentElement);

          groupContainerDiv.appendChild(tile);
        });

        groupContainer.appendChild(groupContainerDiv);
      });
    } else {
      const emptyDataMessage = document.createElement("p");
      emptyDataMessage.textContent = "No data available.";
      groupContainer.appendChild(emptyDataMessage);
    }
  })
  .catch(console.error);
