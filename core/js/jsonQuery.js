function getName(arg) {
  if (typeof arg === 'string' || typeof arg === 'number') {
    return Promise.resolve(arg);
  } else if (typeof arg === 'object' && arg !== null && 'link' in arg && 'path' in arg) {
    const { link, path } = arg;
    return getValueFromURL(link, path);
  } else {
    return Promise.reject(new Error('Invalid argument type or missing required fields.'));
  }
}
  
// Sample function for fetching value from URL and response path (as in the previous examples)
function getValueFromURL(url, responsePath) {
  return fetch(url)
    .then(response => response.json())
    .then(jsonData => {
      const keys = responsePath.split('.');
      let value = jsonData;
      for (const key of keys) {
        if (!value.hasOwnProperty(key)) {
          throw new Error(`Invalid path: ${responsePath}`);
        }
        value = value[key];
      }
      return value;
    })
    .catch(error => {
      console.error('Error fetching data:', error);
      return null;
    });
}
// Export the createContextMenu function
export { getName };
export { getValueFromURL };