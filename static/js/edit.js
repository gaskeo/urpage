function changeTab(evt, tab, section) {
    let i, sections, tabs, activeSection;

    sections = document.getElementsByClassName("section");

    activeSection = document.getElementById(section)

    for (i = 0; i < sections.length; i++) {
        if (sections[i] !== activeSection) {
            sections[i].style.opacity = "0";
            sections[i].style.display = "none";
        } else {
            activeSection.style.opacity = "1"
            activeSection.style.display = "block"
        }

    }

    tabs = document.getElementsByClassName("tabs-button active-tab");

    for (i = 0; i < tabs.length; i++) {
        tabs[i].className = tabs[i].className.replace(" active-tab", "");
    }

    evt.currentTarget.className += " active-tab";
}

function deleteLink(divId) {
    let elem, index, countElements, i
    elem = document.getElementById(divId)
    index = divId.replace("link-div-", "")
    countElements = document.getElementsByClassName("page-form-link").length

    if (parseInt(index) !== countElements - 1) {
        for (i = parseInt(index) + 1; i < countElements; i++) {
            document.getElementById(`link-div-${i}`).id = `link-div-${(i - 1)}`
            document.getElementById(`link-${i}`).id = `link-${(i - 1)}`

            document.getElementById(`delete-link-${i}`).setAttribute("onclick", `deleteLink('link-div-${(i - 1)}')`)

            document.getElementById(`delete-link-${i}`).id = `delete-link-${(i - 1)}`
        }
    }
    elem.remove()
}

function changePreview() {
    let preview = document.getElementById('img_src');

    let file = document.getElementById('img-request').files[0];

    let reader = new FileReader();

    reader.onloadend = function () {
        preview.src = reader.result;
        return reader.result
    }

    if (file) {
        console.log(reader.readAsDataURL(file));
    }

}


function addLink() {
    let newElementN, form, newDiv, newInput, newButton, afterItem
    newElementN = document.getElementsByClassName("page-form-link").length

    form = document.getElementById("page-form-links")
    afterItem = document.getElementById("add-link-button")

    newDiv = document.createElement("div")
    newDiv.id = `link-div-${newElementN}`
    newDiv.style.display = "inline"
    newDiv.style.verticalAlign = "middle"

    newInput = document.createElement("input")
    newInput.id = `link-${newElementN}`
    newInput.name = "link"
    newInput.className = "page-form-input page-form-link"
    newInput.type = "text"

    newButton = document.createElement("button")
    newButton.id = `delete-link-${newElementN}`
    newButton.className = "delete-link-button"
    newButton.setAttribute("onclick", `deleteLink('link-div-${newElementN}')`)
    newButton.textContent = "-"
    newButton.type = "button"

    newDiv.appendChild(newInput)
    newDiv.appendChild(newButton)

    form.insertBefore(newDiv, afterItem)
}

function sendMain() {
    let status

    let data = new FormData();

    let photo = document.getElementById("img-request").files[0];

    data.append("id", document.getElementById("id").value)
    data.append("username", document.getElementById("username").value)
    data.append("image", photo)


    fetch("/do/edit_main", {method: 'post', body: data}).then(function (r) {
        status = r.status

        if (status === 200) {
            alert("Успешно")
        } else {
            alert("что-то пошло не так...")
        }
    })

    return false
}

function sendLinks() {
    return true
}