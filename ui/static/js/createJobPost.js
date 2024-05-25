// edit req/roles list button text content
const createJobPostFormReqList = document.getElementById('createJobPostFormReqList')
const editReqList = document.getElementById('editReqList')
const createJobPostFormRoleList = document.getElementById('createJobPostFormRoleList')
const editRoleList = document.getElementById('editRoleList')

const getListItemEditContent = (content, isEditing) => `
  <textarea maxlength="200" disabled oninput='this.style.height = "";this.style.height = this.scrollHeight + 3 + "px"'>${content}</textarea>
  <div class="req-list__editing-item--buttons">
    <div>
    ${isEditing
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

  return listElement
}

// console.log(createJobPostFormReqList.children)
// console.log(createJobPostFormReqList.children.length)

if (createJobPostFormReqList) {
  if (createJobPostFormReqList.children.length > 0) {
    editReqList.textContent = 'Edit'
  } else {
    editReqList.textContent = 'Add a requirement'
  }
}

if (createJobPostFormRoleList) {
  if (createJobPostFormRoleList.children.length > 0) {
    editRoleList.textContent = 'Edit'
  } else {
    editRoleList.textContent = 'Add a task'
  }
}

// add behavior for edit req list button
if (editReqList) {
  editReqList.addEventListener('click', function (e) {
    //~~~   enter editing state
    console.log(editReqList.dataset.isEditing);
    if (editReqList.dataset.isEditing === 'false') {
      console.log(
        'amount of elements in req list',
        createJobPostFormReqList.children.length
      )
      // const amountOfCurrentItems = createJobPostFormReqList.children.length
      if (createJobPostFormReqList.children.length === 0) {
        createJobPostFormReqList.appendChild(document.createElement('li'))
        createJobPostFormReqList.children[0].innerHTML = getListItemEditContent('', true)
        createJobPostFormReqList.children[0].classList.add('req-list__editing-item')
        createJobPostFormReqList.children[0].children[0].focus()
      } else {
        // get list items text content
        console.log('get list items text content')

        const reqItemsText = []
        const newReqItems = []

        for (let i = 0; i < createJobPostFormReqList.children.length; i++) {
          reqItemsText.push(createJobPostFormReqList.children[i].textContent.trim())

          const newElement = createListElementEditing(
            createJobPostFormReqList.children[i].textContent.trim()
          )
          newReqItems.push(
            newElement
          )

          // createJobPostFormReqList.appendChild(
          //   newElement
          // )
        }
        console.log(reqItemsText);
        console.log(newReqItems);
        //delete child elements
        let lastListElement = createJobPostFormReqList.lastElementChild
        while (lastListElement) {
          createJobPostFormReqList.removeChild(lastListElement)
          lastListElement = createJobPostFormReqList.lastElementChild
        }

        createJobPostFormReqList.append(
          ...newReqItems
        )

        // modify textarea height
        for (let i = 0; i < createJobPostFormReqList.children.length; i++) {
          // target the textarea inside the li element and re-set its height
          createJobPostFormReqList.children[i].children[0].style.height =
            createJobPostFormReqList.children[i].children[0].scrollHeight + 3 + 'px'

            
          // append listeners for edit btn
          createJobPostFormReqList.children[i].children[1].children[0].addEventListener('click', function (e) {
            
            const icon = e.target
            const iconDivParent = e.target.parentElement
            
            if (icon.closest('.req-list__editing-item--buttons').previousElementSibling.disabled === true) {

              console.log('click on', icon);
              console.log('parent el', iconDivParent);
              console.log('last el of parent el', iconDivParent.lastElementChild);
              console.log('this will focus the textarea an allow editing');
              
              console.log(icon.closest('.req-list__editing-item--buttons'))
              console.log(icon.closest('.req-list__editing-item--buttons').previousElementSibling)
              // enable the textarea and focus it
              icon.closest('.req-list__editing-item--buttons').previousElementSibling.disabled = false
              icon.closest('.req-list__editing-item--buttons').previousElementSibling.focus()
              
              
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
              icon.closest('.req-list__editing-item--buttons').previousElementSibling.disabled = true
              // icon.closest('.req-list__editing-item--buttons').previousElementSibling.focus()
              
              
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
          })

          // append listeners for delete btn
          createJobPostFormReqList.children[i].children[1].children[1].addEventListener('click', function (e) {
            const icon = e.target
            const iconDivParent = e.target.parentElement
            const liElementToRemove = iconDivParent.closest(
              '.req-list__editing-item'
            )
            console.log(icon);
            console.log(iconDivParent);

            liElementToRemove.remove()            
          })
        }
      }

      editReqList.dataset.isEditing = 'true'
      editReqList.textContent = 'Done'
      
    } else { // exit editing state
      editReqList.dataset.isEditing = 'false'
      editReqList.textContent = 'Edit'
      
    }

  })
} 