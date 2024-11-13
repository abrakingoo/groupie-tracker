document.addEventListener("DOMContentLoaded", async () => {
  const searchInput = document.getElementById("search-term");
  const suggestionsBox = document.getElementById("suggestions");

  let searchTerm = searchInput.value.toLowerCase().trim();

  function displaySuggestions(suggestions) {
    suggestionsBox.innerHTML = "";
    if (suggestions.length === 0) {
      suggestionsBox.style.display = "none";
      return;
    }
    suggestionsBox.style.display = "block";
    suggestions.forEach((item) => {
      const itemDiv = document.createElement("div");
      itemDiv.textContent = item;
      itemDiv.style.cursor = "pointer";
      itemDiv.className = "member-suggestion";
      itemDiv.tabIndex = 0;
      itemDiv.addEventListener("click", () => {
        searchInput.value = item;
        searchInput.focus();
        suggestionsBox.innerHTML = "";
        suggestionsBox.style.display = "none";
      });
      itemDiv.addEventListener("keydown", (e) => {
        if (e.key === "Enter") {
          itemDiv.click();
        }
      });
      suggestionsBox.appendChild(itemDiv);
    });
  }

  // Hide suggestions box if search term is empty
  if (searchTerm.length === 0) {
    suggestionsBox.style.display = "none";
    suggestionsBox.innerHTML = "";
  }

  const proxyUrl = "https://thingproxy.freeboard.io/fetch/";
  const targetUrl = "https://groupietrackers.herokuapp.com/api/locations";
  let resLocations = [];

  async function fetchLocations() {
    try {
      const res = await fetch(proxyUrl + targetUrl);
      if (!res.ok) {
        throw new Error(`HTTP error! Status: ${res.status}`);
      }
      const data = await res.json();
      resLocations = data.index.flatMap((item) => item.locations);
    } catch (error) {
      console.error("Error fetching locations:", error);
    }
  }

  await fetchLocations();

  const uniq = arr => [...new Set(arr)]

  fetch("http://localhost:8080/bands")
    .then((res) => res.json())
    .then((data) => {
      searchInput.addEventListener("keyup", () => {
        searchTerm = searchInput.value.toLowerCase().trim();
        if (searchTerm.length === 0) {
          suggestionsBox.style.display = "none";
          suggestionsBox.innerHTML = "";
          return;
        }

        const bandMatches = data
          .filter((band) => band.name.toLowerCase().includes(searchTerm))
          .map((band) => 
            `${band.name} - artist/band`,
          );

        const matchingMembers = data.flatMap((band) =>
          band.members
            .filter((member) => member.toLowerCase().includes(searchTerm))
            .map((member) => `${member} - member`)
        );

        const creationDates = data
          .filter((band) => String(band.creationdate).includes(searchTerm))
          .map((band) => 
            `${band.creationdate} - creationdate`,
          );

        const firstAlbums = data
          .filter((band) => String(band.firstAlbum).includes(searchTerm))
          .map((band) => 
            `${band.firstAlbum} - firstAlbum`,
          );

        const locationMatches = resLocations
          .filter((location) => location.toLowerCase().includes(searchTerm))
          .map((location) => 
            `${location} - locations`,
          );

        const band = uniq(bandMatches)
        const members = uniq(matchingMembers)
        const dates = uniq(creationDates)
        const album = uniq(firstAlbums)
        const location = uniq(locationMatches)


        const allMatches = [
          ...band,
          ...members,
          ...dates,
          ...album,
          ...location,
        ];

        displaySuggestions(allMatches);
      });
    });
});
