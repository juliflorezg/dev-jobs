
  function loadJSForManageJP(){
    console.log('afdfds afsad');
    const modal = document.querySelector('#manageJobPostsModal')
    const deleteBtns = document.querySelectorAll('.child-item--buttons-delete')
    const closeModal = document.querySelector('.close-modal-btn')
    const manageJobPostsModalForm = document.getElementById('manageJobPostsModalForm')
    const jpIdToDelete = document.getElementById('jpIdToDelete')
    const modalContentJobPostPosition = document.getElementById('modalContentJobPostPosition')

    deleteBtns.forEach(function (btn) {
      btn.addEventListener('click', () => {
        modal.showModal()
        jpIdToDelete.value = btn.dataset.jpid
        console.log(btn.dataset);
        modalContentJobPostPosition.textContent = btn.dataset.jpPosition
      })
    })

    modal.addEventListener('click', function (e) {
      console.log(e.target.nodeName);
      if (e.target.nodeName === "DIALOG") {
        modal.close()

      }
    })

    closeModal.addEventListener('click', function (e) {
      modal.setAttribute('closing', '')

      modal.addEventListener('animationend', function () {
        modal.removeAttribute('closing')
        modal.close()

      }, { once: true })
    })

    manageJobPostsModalForm.addEventListener('submit', function (e) {
      e.preventDefault()

      console.log(e.target.jpid);

      modal.setAttribute('closing', '')

      modal.addEventListener('animationend', function () {
        modal.removeAttribute('closing')
        modal.close()

      }, { once: true })

      // construct a FormData object, which fires the formdata event
      const formData = new FormData(e.target);
      const jpid = formData.get("jpid")
      console.log({ jpid });

      setTimeout(() => {
        // e.target.submit()

        fetch(e.target.action, {
          method: 'DELETE',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            jpid: jpid
          })
        }).then(res => {
          console.log(res);
          if (res.ok) {
            return res.text()
          }
        }).then(html => {
          console.log(html);

          if (html) {
            window.scrollTo({
              top: 0,
              behavior: 'smooth',
            })
            document.documentElement.innerHTML = html
            // window.history.pushState({}, "My account", "/user/account")

            if (document.title.includes("Manage my JobPosts")) {

              // Dispatch a custom event after replacing the HTML
              const event = new CustomEvent('htmlContentReplaced');
              document.dispatchEvent(event);
            }
          }

        }).catch(err => {
          console.log("error trying to delete a publication:", err);
        })

      }, 1000);

    })
  }



if (document.title.includes("Manage my JobPosts")) {
  document.addEventListener('DOMContentLoaded', function () {
    loadJSForManageJP()
  })

  document.addEventListener('htmlContentReplaced', function () {
    console.log('content has been replaced ');
    loadJSForManageJP()
  })
}