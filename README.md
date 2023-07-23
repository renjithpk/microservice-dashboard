# Microservice Dashboard

The Microservice Dashboard is a web application that provides a centralized view of various microservices within your architecture. It allows you to monitor the health and status of individual microservices, track their performance metrics, and manage their configurations. 

## Configuration
The data.yaml file is used to configure the microservices displayed on the dashboard. It is a YAML file with the following structure:

```yaml
groups:
  - name: Group Name
    color: "#RRGGBB"
    items:
      - name: Microservice 1
        menu:
          - name: Menu Item 1
            link: http://link/microservice1
          - name: Menu Item 2
            link: http://link/microservice1
        buttons:
          - name: "Button 1"
            color: "#RRGGBB"
            items:
              - name: Submenu Item 1
                link: https://somelink/submenu1
              - name: Submenu Item 2
                link: https://somelink/submenu2
        backends: ["BackendMicroservice1", "BackendMicroservice2"]
      - name: Microservice 2
        menu:
          - name: Menu Item 1
            link: http://link/microservice2
          - name: Menu Item 2
            link: http://link/microservice2
        backends: ["BackendMicroservice3"]
  # Add other groups and microservices here
```
Configuration Elements:
groups: The top-level key representing a list of groups. Each group contains a set of related microservices.

`name`: The name of the group. It represents a category of microservices (e.g., "Booking," "Payment," etc.).

`color`: The color used to represent the group in the dashboard. It should be in the format "#RRGGBB," where RR, GG, and BB are hexadecimal values representing the red, green, and blue components of the color, respectively.

`items`: A list of microservices within the group. Each item represents a single microservice.

`name`: The name of the microservice.

`menu`: A list of menu items for the microservice. Each menu item consists of a name and a link representing the label and the URL link, respectively.

`buttons` (Optional): A list of clickable buttons for the microservice. Each button contains a name, color, and items. The items field is a list of submenu items for the button, similar to the menu items.

`backends` (Optional): A list of backend microservices connected to the current microservice. This field helps in establishing relationships between microservices for highlighting and disabling unrelated items.

You can customize the data.yaml file to add or remove groups and microservices, adjust their colors, define menu and submenu items, and establish connections between microservices using the backends field. Ensure that the YAML syntax is correct, and each element is correctly nested within the appropriate structure.