const createJobPostForm = document.getElementById('createJobPostForm')
const createJobPostFormReqList = document.getElementById('createJobPostFormReqList')
const editReqListBtn = document.getElementById('editReqList')
const createJobPostFormRoleList = document.getElementById('createJobPostFormRoleList')
const editRoleListBtn = document.getElementById('editRoleList')

const getListItemEditContent = (content, isEditing) => `
  <textarea maxlength="200" disabled oninput='this.style.height = "";this.style.height = this.scrollHeight + 3 + "px"'>${content}</textarea>
  <div class="req-list__editing-item--buttons">
    <div>
    ${isEditing   // check mark
    ? `<svg
          xmlns="http://www.w3.org/2000/svg"
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          stroke="#ffffff"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="icon icon-tabler icons-tabler-outline icon-tabler-check"
        >
          <path stroke="none" d="M0 0h24v24H0z" fill="none" />
          <path d="M5 12l5 5l10 -10" />x
        </svg>`
    // pencil 
    : `<svg          
          xmlns="http://www.w3.org/2000/svg"
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          stroke="#ffffff"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="icon icon-tabler icons-tabler-outline icon-tabler-pencil"
        >
          <path stroke="none" d="M0 0h24v24H0z" fill="none" />
          <path d="M4 20h4l10.5 -10.5a2.828 2.828 0 1 0 -4 -4l-10.5 10.5v4" />
          <path d="M13.5 6.5l4 4" />
        </svg>`
  }
    </div>
    <div>
      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
        stroke="#ffffff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
        class="icon icon-tabler icons-tabler-outline icon-tabler-trash">
        <path stroke="none" d="M0 0h24v24H0z" fill="none" />
        <path d="M4 7l16 0" />
        <path d="M10 11l0 6" />
        <path d="M14 11l0 6" />
        <path d="M5 7l1 12a2 2 0 0 0 2 2h8a2 2 0 0 0 2 -2l1 -12" />
        <path d="M9 7v-3a1 1 0 0 1 1 -1h4a1 1 0 0 1 1 1v3" />
      </svg>
    </div>
  </div>
`

function createListElementEditing(content) {
  const listElement = document.createElement('li')
  listElement.innerHTML = getListItemEditContent(content, false)
  listElement.classList.add('req-list__editing-item')

  listElement.children[0].addEventListener('input', function (e) {
    console.log(e.target.value);
    e.target.textContent = e.target.value
  })

  return listElement
}

function listenerForEditBtn(target) {
  const icon = target
  const iconDivParent = target.parentElement
  console.log('click on', icon)
  console.log('parent el', iconDivParent)

  if (
    icon.closest('.req-list__editing-item--buttons').previousElementSibling
      .disabled === true
  ) {
    console.log('click on', icon)
    console.log('parent el', iconDivParent)
    console.log('last el of parent el', iconDivParent.lastElementChild)
    console.log('this will focus the textarea an allow editing')

    console.log(icon.closest('.req-list__editing-item--buttons'))
    console.log(
      icon.closest('.req-list__editing-item--buttons').previousElementSibling
    )
    // enable the textarea and focus it
    icon.closest(
      '.req-list__editing-item--buttons'
    ).previousElementSibling.disabled = false
    icon
      .closest('.req-list__editing-item--buttons')
      .previousElementSibling.focus()

    //change the icon
    let iconSVG = iconDivParent.lastElementChild
    while (iconSVG) {
      iconDivParent.removeChild(iconSVG)
      iconSVG = iconDivParent.lastElementChild
    }
    iconDivParent.innerHTML = `<svg
                xmlns="http://www.w3.org/2000/svg"
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                stroke="#ffffff"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="icon icon-tabler icons-tabler-outline icon-tabler-check"
                >
                <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                <path d="M5 12l5 5l10 -10" />
                </svg>`
  } else {
    icon.closest(
      '.req-list__editing-item--buttons'
    ).previousElementSibling.disabled = true

    //set textarea.content
    console.log(
      icon.closest('.req-list__editing-item--buttons').previousElementSibling
        .textContent
    )

    //change the icon
    let iconSVG = iconDivParent.lastElementChild
    while (iconSVG) {
      iconDivParent.removeChild(iconSVG)
      iconSVG = iconDivParent.lastElementChild
    }
    iconDivParent.innerHTML = `<svg
                xmlns="http://www.w3.org/2000/svg"
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                stroke="#ffffff"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="icon icon-tabler icons-tabler-outline icon-tabler-pencil"
                >
                <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                <path d="M4 20h4l10.5 -10.5a2.828 2.828 0 1 0 -4 -4l-10.5 10.5v4" />
                <path d="M13.5 6.5l4 4" />
                </svg>`
  }
}

function listenerForDeleteBtn(target) {
  const icon = target
  const iconDivParent = target.parentElement
  const liElementToRemove = iconDivParent.closest('.req-list__editing-item')
  console.log(icon)
  console.log(iconDivParent)

  liElementToRemove.remove()
}

function manageEditButtonAction(btn, list) {
  console.log(btn.dataset.isEditing)
  if (btn.dataset.isEditing === 'false') {
    //~~~   enter editing state

    //> add another button before edit btn that says 'add an element'
    {
      const newBtnForAddLI = document.createElement('BUTTON')
      newBtnForAddLI.textContent = 'Add an element'
      newBtnForAddLI.classList.add('list-buttons__normal-btn')
      btn.parentElement.insertBefore(newBtnForAddLI, btn)

      newBtnForAddLI.addEventListener('click', function (e) {
        e.preventDefault()
        list.appendChild(createListElementEditing(''))

        let currentAmountOfLis = list.children.length
        console.log(currentAmountOfLis)

        list.children[
          currentAmountOfLis - 1
        ].children[0].disabled = false
        list.children[
          currentAmountOfLis - 1
        ].children[0].focus()

        //> add listeners to buttons
        list.children[
          currentAmountOfLis - 1
        ].children[1].children[0].addEventListener('click', function (e) {
          listenerForEditBtn(e.target)
        })

        list.children[
          currentAmountOfLis - 1
        ].children[1].children[1].addEventListener('click', function (e) {
          listenerForDeleteBtn(e.target)
        })
      })
    }

    console.log(
      'amount of elements in req list',
      list.children.length
    )
    if (list.children.length === 0) {
      // add a single list item ready for enter description
      list.appendChild(document.createElement('li'))
      list.children[0].classList.add(
        'req-list__editing-item'
      )
      list.children[0].innerHTML = getListItemEditContent(
        '',
        true
      )
      list.children[0].children[0].disabled = false
      list.children[0].children[0].focus()
    } else {
      // get list items text content
      console.log('get list items text content')

      // const reqItemsText = []
      const newReqItems = []

      for (let i = 0; i < list.children.length; i++) {
        // reqItemsText.push(list.children[i].textContent.trim())

        const newElement = createListElementEditing(
          list.children[i].textContent.trim()
        )
        newReqItems.push(newElement)

        // list.appendChild(
        //   newElement
        // )
      }
      // console.log(reqItemsText);
      console.log(newReqItems)
      //delete child elements
      let lastListElement = list.lastElementChild
      while (lastListElement) {
        list.removeChild(lastListElement)
        lastListElement = list.lastElementChild
      }

      list.append(...newReqItems)

      for (let i = 0; i < list.children.length; i++) {
        const listElement = list.children[i]

        // modify textarea height
        // target the textarea inside the li element and re-set its height
        listElement.children[0].style.height =
          listElement.children[0].scrollHeight + 3 + 'px'

        // append listeners for edit btn
        listElement.children[1].children[0].addEventListener(
          'click',
          function (e) {
            listenerForEditBtn(e.target)
          }
        )

        // append listeners for delete btn
        listElement.children[1].children[1].addEventListener(
          'click',
          function (e) {
            listenerForDeleteBtn(e.target)
          }
        )
      }
    }

    btn.dataset.isEditing = 'true'
    btn.textContent = 'Done'
  } else {
    // exit editing state

    btn.dataset.isEditing = 'false'
    btn.textContent = 'Edit'
    btn.previousElementSibling.remove()

    // loop over the list and get the text in each text area (only those that have a text other than "")
    const newReqItems = []

    for (let i = 0; i < list.children.length; i++) {
      const li = list.children[i]

      // console.log(li);
      // console.log(li.children[0].textContent);

      if (li.children[0].textContent.trim().length > 0) {
        const newElement = document.createElement('LI')
        newElement.textContent = li.children[0].textContent.trim()
        newReqItems.push(newElement)
      } else {
        continue
      }
    }
    // console.log(reqItemsText);
    console.log(newReqItems)
    //delete child elements
    let lastListElement = list.lastElementChild
    while (lastListElement) {
      list.removeChild(lastListElement)
      lastListElement = list.lastElementChild
    }

    list.append(...newReqItems)
  }
}

if (createJobPostFormReqList) {
  if (createJobPostFormReqList.children.length > 0) {
    editReqListBtn.textContent = 'Edit'
  } else {
    editReqListBtn.textContent = 'Add a requirement'
  }
}

if (createJobPostFormRoleList) {
  if (createJobPostFormRoleList.children.length > 0) {
    editRoleListBtn.textContent = 'Edit'
  } else {
    editRoleListBtn.textContent = 'Add a task'
  }
}

// add behavior for edit req list button
if (editReqListBtn) {
  editReqListBtn.addEventListener('click', function (e) {
    console.log(editReqListBtn.dataset.isEditing);

    manageEditButtonAction(editReqListBtn, createJobPostFormReqList)
  })
} 

if (editRoleListBtn) {
  editRoleListBtn.addEventListener('click', function (e) {
    console.log(editRoleListBtn.dataset.isEditing);

    manageEditButtonAction(editRoleListBtn, createJobPostFormRoleList)
  })
} 


// manage form submission
createJobPostForm.addEventListener('submit', function(e){
  e.preventDefault()

  if (editReqListBtn.dataset.isEditing === "true"){
    console.log('still editing reqs, please confirm you requirements and try again');

    editReqListBtn.closest('.create-jobpost-container__form-req-list-container').classList.add('border-error-xs')
  }


  // fetch('/account/create/jobpost', {
  //   body: JSON.stringify(jobPostData)
  // })
})