function changeTab(evt, tab, section) {
    let i, sections, tabs;

    sections = document.getElementsByClassName("section");

    for (i = 0; i < sections.length; i++) {
        sections[i].style.opacity = "0";
        sections[i].style.display = "none";
    }

    tabs = document.getElementsByClassName("tabs-button active-tab");

    for (i = 0; i < tabs.length; i++) {
        tabs[i].className = tabs[i].className.replace(" active-tab", "");
    }

    evt.currentTarget.className += " active-tab";
    document.getElementById(section).style.display = "block";

    document.getElementById(section).style.opacity = "1";

}

window.onload = function () {
    document.getElementById('img_click').addEventListener('click', function (e) {
        let input = document.getElementById("img");
        input.type = 'file';
        input.accept = ".jpg, .jpeg, .png";
        input.onchange = e => {
            let file = e.target.files[0];
            document.getElementById("img_src").src = window.URL.createObjectURL(file);
            form.image = input
        };
        input.click();
    });
};

function deleteLink(divId) {
    let elem, index, elemsN, i
    elem = document.getElementById(divId)
    index = divId.replace("link-div-", "")
    elemsN = document.getElementsByClassName("page-form-link").length

    if (parseInt(index) !== elemsN - 1) {
        for (i = parseInt(index) + 1; i < elemsN; i++) {
            document.getElementById("link-div-" + i).id = "link-div-" + (i - 1)
            document.getElementById("link-" + i).id = "link-" + (i - 1)

            document.getElementById("delete-link-" + i).setAttribute("onclick", "deleteLink('link-div-" + (i - 1) + "')")

            document.getElementById("delete-link-" + i).id = "delete-link-" + (i - 1)
        }
    }
    elem.remove()
}


function addLink() {
    let elemsN, form, newDiv, newInput, newButton, afterItem, newP
    elemsN = document.getElementsByClassName("page-form-link").length

    form = document.getElementById("page-form-links")
    afterItem  = document.getElementById("add-link-button")

    newDiv = document.createElement("div")
    newDiv.id = "link-div-" + elemsN
    newDiv.style.display = "inline"
    newDiv.style.verticalAlign = "middle"

    newInput = document.createElement("input")
    newInput.id = "link-" + elemsN
    newInput.name = "link"
    newInput.className = "page-form-input page-form-link"
    newInput.type = "text"

    newButton = document.createElement("button")
    newButton.id = "delete-link-" + elemsN
    newButton.className = "delete-link-button"
    newButton.setAttribute("onclick", "deleteLink('link-div-" + elemsN + "')")
    newButton.textContent = "-"
    newButton.type = "button"

    newDiv.appendChild(newInput)
    newDiv.appendChild(newButton)

    form.insertBefore(newDiv, afterItem)
}