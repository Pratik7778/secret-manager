import React, { useState, useEffect } from 'react';
import {Button, Modal} from 'react-bootstrap';

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faTrash } from '@fortawesome/free-solid-svg-icons';
import Cookies from 'js-cookie';

const DeleteSecretPage = ({secretName, refreshFlag, setRefreshFlag}) => {
    const [message, setMessage] = useState('')
    const [showModal, setShowModal] = useState(false)

    const handleShow = () => setShowModal(true)
    const handleClose = () => setShowModal(false)
    const deleteSecret = async () => {
        const token = Cookies.get("token")
        if(!token) {
            setMessage({type: "error", message: "Unauthorized: No token found."})
            return
        }

        try {
            // const API_URL = await import.meta.env.VITE_API_URL
            const API_URL = localStorage.getItem("API_URL") || "http://localhost:8080"
            const response = await fetch(API_URL + `/api/v1/secrets/${secretName}`, {
                method: "DELETE",
                headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
                },
            })

            const result = await response.json()

            if(response.status === 200) {
                setMessage({type: "success", message: result.message})
                setRefreshFlag(!refreshFlag)
                
                console.log(result)

            } else if(response.status === 400) {
                setMessage({type: "error", message: result.error})
                console.log(result)

            } else if(response.status === 401) {
                setMessage({type: "error", message: result}) 
                console.log(result)
            }

        } catch (err) {
            console.log(err)
        }
    }

    return (
        <>
        <td>
            <FontAwesomeIcon 
                icon={faTrash} 
                size="1x" 
                onClick={handleShow} 
                style={{ cursor: 'pointer', marginLeft: "-460px"}} 
            />

            <Modal show={showModal} onHide={handleClose}>
                <Modal.Header closeButton>
                <Modal.Title>Confirm Action</Modal.Title>
                </Modal.Header>
                <Modal.Body>
                    <p>Are you sure you want to delete this item?</p>
                </Modal.Body>
                <Modal.Footer>
                    <Button variant="danger" onClick={deleteSecret}>
                        Confirm
                    </Button>
                </Modal.Footer>
            </Modal>
        </td>
        </>
    )
}

export default DeleteSecretPage;
