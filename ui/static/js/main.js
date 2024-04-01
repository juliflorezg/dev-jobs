formMobileMenu = document.querySelector('.form-input-container-mobile')
openFormMenuIcon = document.querySelector(
  `.form-filter-container-search-mobile > button[type='button']`
)
closeFormIcon = document.querySelector('.close-icon-container')

formMobileMenu.addEventListener('click', function (e) {
  if (e.target === formMobileMenu) {
    console.log('touched on mobile menu')
    formMobileMenu.classList.add('hidden')
  }
})

closeFormIcon.addEventListener('click', function (e) {
  console.log('click on' + e.target)
  formMobileMenu.classList.add('hidden')
})

openFormMenuIcon.addEventListener('click', function (e) {
  formMobileMenu.classList.remove('hidden')
})

document.addEventListener('DOMContentLoaded', function (e) {
  console.log(window.innerWidth)

  if (window.innerWidth >= 768) {
    formMobileMenu.classList.remove('hidden')
  }
})
