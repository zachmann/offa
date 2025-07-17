const dropdown = document.querySelector('.dropdown');
const toggle = dropdown.querySelector('.dropdown-toggle');
const menu = dropdown.querySelector('.dropdown-menu');
const input = menu.querySelector('input');
const options = menu.querySelectorAll('.option');
const formInput = document.getElementById('issuer');
const form = document.getElementById('form');

toggle.addEventListener('click', () => {
    menu.style.display = menu.style.display === 'block' ? 'none' : 'block';
    input.focus();
});

options.forEach(option => {
    option.addEventListener('click', () => {
        toggle.innerHTML = option.innerHTML;
        const value = option.getAttribute('data-value');
        formInput.value = value;
        menu.style.display = 'none';
        form.submit();
    });
});

window.addEventListener('click', (e) => {
    if (!dropdown.contains(e.target)) {
        menu.style.display = 'none';
    }
});

input.addEventListener('input', () => {
    const filter = input.value.toLowerCase();
    options.forEach(option => {
        const text = option.textContent.toLowerCase();
        const tags = option.getAttribute('data-tags').toLowerCase();
        if (text.includes(filter) || tags.includes(filter)) {
            option.style.display = 'flex';
        } else {
            option.style.display = 'none';
        }
    });
});