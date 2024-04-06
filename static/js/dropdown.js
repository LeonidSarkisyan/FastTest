const plus = document.getElementById("plus")
const dropDown = document.getElementById("myDropdown")

plus.addEventListener("mouseenter", toggleDropdown)

function toggleDropdown() {
    let dropdown = document.getElementById("myDropdown");
    dropdown.classList.toggle("show");
}

let timeout

plus.addEventListener("mouseenter", function () {
    clearTimeout(timeout)
})

plus.addEventListener("mouseleave", function(event) {
    timeout = setTimeout(
        function () {
            let dropdowns = document.getElementsByClassName("dropdown-content");
            for (let i = 0; i < dropdowns.length; i++) {
                let openDropdown = dropdowns[i];
                if (openDropdown.classList.contains('show')) {
                    openDropdown.classList.remove('show');
                }
            }
        }
    , 250)
})

dropDown.addEventListener("mouseenter", function () {
    clearTimeout(timeout)
})

dropDown.addEventListener("mouseleave", function () {
    timeout = setTimeout(function () {
        let dropdowns = document.getElementsByClassName("dropdown-content");
        for (let i = 0; i < dropdowns.length; i++) {
            let openDropdown = dropdowns[i];
            if (openDropdown.classList.contains('show')) {
                openDropdown.classList.remove('show');
            }
        }}, 250)
})

dropDown.onclick = function () {
    let dropdowns = document.getElementsByClassName("dropdown-content");
    for (let i = 0; i < dropdowns.length; i++) {
        let openDropdown = dropdowns[i];
        if (openDropdown.classList.contains('show')) {
            openDropdown.classList.remove('show');
        }
    }
    setTimeout(function () {
        let dropdowns = document.getElementsByClassName("dropdown-content");
        for (let i = 0; i < dropdowns.length; i++) {
            let openDropdown = dropdowns[i];
            if (openDropdown.classList.contains('show')) {
                openDropdown.classList.remove('show');
            }
        }
    }, 0)
}