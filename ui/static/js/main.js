formMobileMenu = document.querySelector('.form-input-container-mobile')
openFormMenuIcon = document.querySelector(
  `.form-filter-container-search-mobile > button[type='button']`
)
closeFormIcon = document.querySelector('.close-icon-container')

headerSignInBtn = document.querySelector('#headerSignInBtn')
signInBtnMenu = document.querySelector(
  '.header-access-container__sign-in-users-list'
)
headerSignUpBtn = document.querySelector('#headerSignUpBtn')
signUpBtnMenu = document.querySelector(
  '.header-access-container__sign-up-users-list'
)

companySignUpForm = document.querySelector('#companySignUpForm')

userSignupFormInfoMessage = document.querySelector('#userSignupFormInfoMessage')
userSignupFormInfoIcon = document.querySelector(
  '.user-signup-form__info-container--icon-container'
)

userSignupInfoList = document.querySelector(
  '.user-signup-form__info-container--info-list'
)
signupFormInfoCloseIcon = document.querySelector(
  '.info-container__info-list--close-icon-container'
)

if (formMobileMenu) {
  formMobileMenu.addEventListener('click', function (e) {
    if (e.target === formMobileMenu) {
      console.log('touched on mobile menu')
      formMobileMenu.classList.add('hidden')
    }
  })
}

if (closeFormIcon) {
  closeFormIcon.addEventListener('click', function (e) {
    console.log('click on' + e.target)
    formMobileMenu.classList.add('hidden')
  })
}

if (openFormMenuIcon) {
  openFormMenuIcon.addEventListener('click', function (e) {
    formMobileMenu.classList.remove('hidden')
  })
}

if (companySignUpForm) {
}
document.addEventListener('DOMContentLoaded', function (e) {
  console.log(window.innerWidth)

  if (window.innerWidth >= 768 && formMobileMenu) {
    formMobileMenu.classList.remove('hidden')
  }
})

if (headerSignInBtn) {
  headerSignInBtn.addEventListener('click', function (e) {
    signInBtnMenu.classList.toggle('hidden')
    signUpBtnMenu.classList.add('hidden')
  })
}

if (headerSignUpBtn) {
  headerSignUpBtn.addEventListener('click', function (e) {
    signUpBtnMenu.classList.toggle('hidden')
    signInBtnMenu.classList.add('hidden')
  })
}


if (userSignupFormInfoMessage) {
  userSignupFormInfoMessage.addEventListener('click', function (e) {
    userSignupInfoList.classList.toggle('hidden')
  })
}

if (userSignupFormInfoIcon) {
  userSignupFormInfoIcon.addEventListener('click', function (e) {
    userSignupInfoList.classList.toggle('hidden')
  })
}

if (signupFormInfoCloseIcon) {
  signupFormInfoCloseIcon.addEventListener('click', function (e) {
    console.log(e.target.tagName)
    userSignupInfoList.classList.toggle('hidden')
  })
}

if (!document.title.includes('Home')) {
  document.querySelector('header').style.height = '15vh'
  document.querySelector('header .wrapper').style.height = '30%'
  document.querySelector('header .wrapper').style.transform = 'translateY(-20%)'
}
