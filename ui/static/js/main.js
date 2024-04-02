formMobileMenu = document.querySelector('.form-input-container-mobile')
openFormMenuIcon = document.querySelector(
  `.form-filter-container-search-mobile > button[type='button']`
)
closeFormIcon = document.querySelector('.close-icon-container')
homeFilterForm = document.querySelector('#homeFilterForm')
homeFilterFormWindowInput = document.querySelector('#homeFilterFormWindowInput')
homeFilterFormMobileMenuInput = document.querySelector(
  '#homeFilterFormMobileMenuInput'
)

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

// form submission
homeFilterForm.addEventListener('submit', function (e) {
  let windowWidth = window.innerWidth
  const formMenuClasses = formMobileMenu.classList  

  homeFilterFormWindowInput.setAttribute('value', windowWidth)
  homeFilterFormMobileMenuInput.setAttribute('value', formMenuClasses.value)
})
