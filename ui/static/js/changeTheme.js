const root = document.documentElement
const toggleBtn = document.querySelector('.toggle-theme-container__button')
let currentTheme = 'light' // ! GLOBAL

function changeTheme(theme) {
  if (theme === 'light') {  
    root.style.setProperty('--bg-color', 'hsl(220, 29%, 10%)')
    root.style.setProperty('--card-color', 'hsl(219, 29%, 14%)')
    root.style.setProperty('--text-titles', 'hsl(0, 0%, 100%')
    root.style.setProperty('--text-paragraph', 'hsl(214, 17%, 51%)')
    toggleBtn.dataset.activeTheme = 'dark'
    
  } else if (theme === 'dark'){
    root.style.setProperty('--bg-color', 'hsl(210, 22%, 96%')
    root.style.setProperty('--card-color', 'hsl(0, 0%, 100%')
    root.style.setProperty('--text-titles', 'hsl(219, 29%, 14%')
    root.style.setProperty('--text-paragraph', 'hsl(210, 29%, 40%)')
    toggleBtn.dataset.activeTheme = 'light'
  }
}

toggleBtn.addEventListener('click', function(e) {
  currentTheme = e.target.dataset.activeTheme
  console.log(e.target);

  console.log(e.target.dataset);

  changeTheme(currentTheme)
})

document.addEventListener('DOMContentLoaded', function(){
  changeTheme(currentTheme)
})

