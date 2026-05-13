function rand(min, max) {
  const minCeiled = Math.ceil(min);
  const maxFloored = Math.floor(max);
  return Math.floor(Math.random() * (maxFloored - minCeiled + 1) + minCeiled);
}

const ichar = " "; // invisible character
const rawlines = document.getElementById("rawdata").innerHTML;
const lines = rawlines.replaceAll("\n", "").split("<br />");

let stretched_lines = [];

for (const line of lines) {
  let stretched_line = "";

  let reading_tag = false;
  for (let char of line) {
    if (char == ">") {
      reading_tag = false;
      stretched_line += char;
      continue;
    }
    if (char == "<" || reading_tag) {
      reading_tag = true;
      stretched_line += char;
      continue;
    }

    stretched_line += ichar.repeat(rand(0, 15));
    stretched_line += char;
  }

  stretched_lines.push(stretched_line);
}

document.getElementById("out").innerHTML = stretched_lines.join("<br>");

let tickint;
if (localStorage.getItem("played") == "true") {
  document.getElementById("out").innerHTML =
    document.getElementById("rawdata").innerHTML;
} else {
  tickint = setInterval(tick, 7);
}

function tick() {
  const lines = document.getElementById("out").innerHTML.split("<br>");
  let outlines = [];

  for (let line of lines) {
    outlines.push(line.replace(ichar, ""));
  }

  document.getElementById("out").innerHTML = outlines.join("<br>");

  if (lines.join("<br>") == outlines.join("<br>")) {
    clearInterval(tickint);
    localStorage.setItem("played", "true");
  }
}

function replay() {
  localStorage.clear();
  location.reload();
  return false;
}
