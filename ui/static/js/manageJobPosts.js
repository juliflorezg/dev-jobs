if (document.title.includes("Manage my JobPosts")) {

  console.log('afdfds afsad');
  const modal = document.querySelector('#manageJobPostsModal')
  const deleteBtns = document.querySelectorAll('.child-item--buttons-delete')
  const closeModal = document.querySelector('.close-modal-btn')
  const manageJobPostsModalForm = document.getElementById('manageJobPostsModalForm')
  const jpIdToDelete = document.getElementById('jpIdToDelete')

  deleteBtns.forEach(function (btn) {
    btn.addEventListener('click', () => {
      modal.showModal()
      jpIdToDelete.value = btn.dataset.jpid
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

  manageJobPostsModalForm.addEventListener('submit', function(e){
    e.preventDefault()

    console.log(e.target.jpid);

    // construct a FormData object, which fires the formdata event
    const formData = new FormData(e.target);
    const jpid = formData.get("jpid")
    console.log({jpid});

    setTimeout(() => {
      // e.target.submit()

      fetch(e.target.action, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          jpid:jpid
        })
      }).then(res =>{
        console.log(res);
        // if (res.ok){
        //   console.log(res);
        // }
      }).catch(err => {
        console.log("error trying to delete a publication:", err);
      })

    }, 1000);

  })

}