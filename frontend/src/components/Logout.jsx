import React, { useState, useEffect } from 'react';
import {Button, Modal} from 'react-bootstrap';

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSignOut, faTrash } from '@fortawesome/free-solid-svg-icons';

import { useNavigate } from 'react-router-dom';
import Cookies from 'js-cookie';

const LogoutPage = () => {
    const [showModal, setShowModal] = useState(false)

    const handleShow = () => setShowModal(true)
    const handleClose = () => setShowModal(false)

    const navigate = useNavigate()

    const logoutHandler = () => {
        Cookies.remove("token")
        navigate("/login")
    }

    return (
        <>
            <div className="d-flex justify-content-end mb-3">
                <FontAwesomeIcon 
                    icon={faSignOut} 
                    size="2x" 
                    onClick={handleShow} 
                    style={{ cursor: 'pointer' }} 
                />
            </div>

            <Modal show={showModal} onHide={handleClose}>
                <Modal.Header closeButton>
                <Modal.Title>Logout</Modal.Title>
                </Modal.Header>
                <Modal.Body>
                    <p>Are you sure you want to Logout?</p>
                </Modal.Body>
                <Modal.Footer>
                    <Button variant="danger" onClick={logoutHandler}>
                        Yes
                    </Button>
                </Modal.Footer>
            </Modal>
        </>
    )
}

export default LogoutPage;
