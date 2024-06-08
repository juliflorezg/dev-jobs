
function loadInitialJSForCreateJP(){

  const createJobPostForm = document.getElementById('createJobPostForm')
  const createJobPostFormReqList = document.getElementById('createJobPostFormReqList')
  const editReqListBtn = document.getElementById('editReqList')
  const createJobPostFormRoleList = document.getElementById('createJobPostFormRoleList')
  const editRoleListBtn = document.getElementById('editRoleList')
  const createJobPostFormErrorMessage = document.getElementById('createJobPostFormErrorMessage')

  const createJobPostFormDescription = document.getElementById('createJobPostFormDescription')
  const createJobPostFormReqDescription = document.getElementById('createJobPostFormReqDescription')
  const createJobPostFormRoleDescription = document.getElementById('createJobPostFormRoleDescription')

  const INITIAL_BTNS_TEXT = {
    reqs: 'Add a requirement',
    roles: 'Add a task'
  }

  const errorLabels = document.querySelector('.create-jobpost-container__form').querySelectorAll('label.form-error');

  console.log(errorLabels);

  errorLabels.forEach(label => {
    // Find the previous sibling input
    const field = label.previousElementSibling;

    // If the previous sibling is an input, add the 'input-error' class
    if (field && (field.tagName.toLowerCase() === 'input' 
    || field.tagName.toLowerCase() === 'textarea' 
    || field.tagName.toLowerCase() === 'select' 
    || field.classList[0] === "create-jobpost-container__form-req-list-container")) {
      field.classList.add('field-error');
    }
  })


  const textAreaListener = (target) => target.textContent = target.value

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

  function createListElementEditing(content, isEditing) {
    const listElement = document.createElement('li')
    listElement.innerHTML = getListItemEditContent(content, isEditing)
    listElement.classList.add('req-list__editing-item')

    listElement.children[0].addEventListener('input', function (e) {
      textAreaListener(e.target)
    })

    return listElement
  }

  function listenerForEditBtn(target) {
    const icon = target
    const iconDivParent = target.parentElement

    if (
      icon.closest('.req-list__editing-item--buttons').previousElementSibling
        .disabled === true
    ) {
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

    liElementToRemove.remove()
  }

  function manageEditButtonAction(btn, list) {

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
          list.appendChild(createListElementEditing('', true))

          let currentAmountOfLis = list.children.length

          list.children[currentAmountOfLis - 1].children[0].disabled = false
          list.children[currentAmountOfLis - 1].children[0].focus()

          //> add listeners to buttons
          list.children[currentAmountOfLis - 1].children[1].children[0].addEventListener('click', function (e) {
            listenerForEditBtn(e.target)
          })

          list.children[currentAmountOfLis - 1].children[1].children[1].addEventListener('click', function (e) {
            listenerForDeleteBtn(e.target)
          })
        })
      }


      if (list.children.length === 0) {
        // add a single list item ready for enter description

        list.appendChild(createListElementEditing('', true))

        list.children[0].children[0].disabled = false
        list.children[0].children[0].focus()

        //> add listeners for check mark and delete btns
        // append listeners for edit btn
        list.children[0].children[1].children[0].addEventListener(
          'click',
          function (e) {
            listenerForEditBtn(e.target)
          }
        )

        // append listeners for delete btn
        list.children[0].children[1].children[1].addEventListener(
          'click',
          function (e) {
            listenerForDeleteBtn(e.target)
          }
        )

      } else {
        // get list items text content

        // const reqItemsText = []
        const newReqItems = []

        for (let i = 0; i < list.children.length; i++) {
          // reqItemsText.push(list.children[i].textContent.trim())

          const newElement = createListElementEditing(
            list.children[i].textContent.trim(), false
          )
          newReqItems.push(newElement)

          // list.appendChild(
          //   newElement
          // )
        }
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
      btn.previousElementSibling.remove()
      list.closest('.create-jobpost-container__form-req-list-container').classList.remove('border-error-xs')

      // loop over the list and get the text in each text area (only those that have a text other than "")
      const newReqItems = []

      for (let i = 0; i < list.children.length; i++) {
        const li = list.children[i]

        if (li.children[0].textContent.trim().length > 0) {
          const newElement = document.createElement('LI')
          newElement.textContent = li.children[0].textContent.trim()
          newReqItems.push(newElement)
        } else {
          continue
        }
      }
      //delete child elements
      let lastListElement = list.lastElementChild
      while (lastListElement) {
        list.removeChild(lastListElement)
        lastListElement = list.lastElementChild
      }

      list.append(...newReqItems)

      if (list.children.length === 0) {

        if (btn.closest('.create-jobpost-container__form-req-list-container').previousElementSibling.textContent.includes('requirements')) {
          btn.textContent = INITIAL_BTNS_TEXT.reqs
        } else {
          btn.textContent = INITIAL_BTNS_TEXT.roles
        }
      } else {
        btn.textContent = 'Edit'
      }
    }
  }

  function getListItemsTexts(listID) {
    list = document.getElementById(listID)

    texts = []
    for (let i = 0; i < list.children.length; i++) {
      texts.push(list.children[i].textContent);
    }
    return texts
  }


  createJobPostFormDescription.addEventListener('input', function (e) { textAreaListener(e.target) })
  createJobPostFormReqDescription.addEventListener('input', function (e) { textAreaListener(e.target) })
  createJobPostFormRoleDescription.addEventListener('input', function (e) { textAreaListener(e.target) })



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
      manageEditButtonAction(editReqListBtn, createJobPostFormReqList)
    })
  }

  if (editRoleListBtn) {
    editRoleListBtn.addEventListener('click', function (e) {
      manageEditButtonAction(editRoleListBtn, createJobPostFormRoleList)
    })
  }


  // manage form submission
  createJobPostForm.addEventListener('submit', function (e) {
    e.preventDefault()

    if (editReqListBtn.dataset.isEditing === "true") {
      editReqListBtn.closest('.create-jobpost-container__form-req-list-container').classList.add('border-error-xs')
      createJobPostFormErrorMessage.classList.toggle('d-none')
      createJobPostFormErrorMessage.textContent = `Please confirm you requirement list and try again.`
      return
    }

    if (editRoleListBtn.dataset.isEditing === "true") {
      editRoleListBtn.closest('.create-jobpost-container__form-req-list-container').classList.add('border-error-xs')
      createJobPostFormErrorMessage.classList.toggle('d-none')
      createJobPostFormErrorMessage.textContent = `Please confirm you role list and try again.`
      return
    }

    jobPostData = {
      position: document.getElementById('createJobpostFormPosition').value,
      description: createJobPostFormDescription.textContent,
      contract: document.getElementById('createJobpostFormContract').value,
      location: document.getElementById('createJobpostFormLocation').value,
      requirements: {
        content: createJobPostFormReqDescription.textContent,
        items: getListItemsTexts('createJobPostFormReqList')
      },
      role: {
        content: createJobPostFormRoleDescription.textContent,
        items: getListItemsTexts('createJobPostFormRoleList')
      },
    }

    const fakeData = {
      // position: 'ios engineer',
      position: 'App & Website Designer',
      description: "We are looking for thoughtful, well-rounded iOS engineer to join our team. We're looking for someone to help build out the foundation of the app and infrastructure. If you are interested in taking part in building an application that millions of people use every day to increase their productivity, this is the perfect opportunity. You will play an important part in our mobile engineering practice, implementing new features, improving performance, and building beautiful user interfaces.",
      contract: "Full Time",
      location: "New Zealand",
      // location: "",
      requirements: {
        content: "**You are an experienced mobile engineer looking to make Pomodoro one of the best mobile experiences out there. You are someone who excels at customer-centric product development and has a passion for working on application architecture and design, and making smooth, delightful experiences. You care deeply about quality, are energized by partnership and collaboration, and you strive to enable others around you to excel.",
        // content: "",
        items: ["**req 1", "req 2"]
        // items: []
      },
      role: {
        content: "**%You will be responsible for building infrastructure and abstractions to help us double our engineering velocity.  You will work at all layers of the stack and closely with partners across engineering, data science, research, product, and design. You will help our codebase stay ahead of the curve of the constantly evolving development ecosystem.",
        items: ["task 1#@!", "task 2"]
        // items: []
      },
      // secretCode: "69420lol"
    }

    const fakeDataStr = JSON.stringify(fakeData)
    console.log(fakeDataStr);


    fetch('/account/create/jobpost', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      // body: JSON.stringify(jobPostData)
      body: (fakeDataStr)
    })
      .then(res => {
        console.log(res)
        if (!res.ok) {
          if (res.status === 422) {
            return res.text()
          }
        }
      })
      .then(html => {
        document.documentElement.innerHTML = html

        // Dispatch a custom event after replacing the HTML
        const event = new CustomEvent('htmlContentReplaced');
        document.dispatchEvent(event);
      })
      .catch(error => console.log(error))
  })


  //"{\"position\":\"ios engineer\",\"description\":\"curl: (60) SSL certificate problem: self-signed certificate. More details here: https://curl.se/docs/sslcerts.html. curl failed to verify the legitimacy of the server and therefore could not\nestablish a secure connection to it. To learn more about this situation and how to fix it, please visit the web page mentioned above.\",\"contract\":\"full_time\",\"location\":\"New Zealand\",\"requirements\":{\"content\":\"dssa dsadfs sdf dfdsf dsf dsf dsaf \",\"items\":[\"req 1\",\"req 2\"]},\"role\":{\"content\":\"jlkjl klkj lkj lkj\",\"items\":[\"task 1\",\"task2\"]}}"

  //

}

document.addEventListener('DOMContentLoaded', function (){
  loadInitialJSForCreateJP()
})

document.addEventListener('htmlContentReplaced', function(){
  console.log('content has been replaced ');
  loadInitialJSForCreateJP()
})