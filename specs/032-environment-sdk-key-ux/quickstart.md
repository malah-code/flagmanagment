# Quickstart Validation: Environment SDK Key UX & Integration Guide

## Prerequisites
- Frontend development server running (`npm run dev` in `frontend/`)
- Backend server running

## Validation Steps

1. **Environment Key Display**:
   - Navigate to Environment Settings under a project.
   - Verify that each environment card displays the Client SDK key with a "Public / Client Key" badge.
   - Click "Copy Key" and verify the key is saved to your clipboard.

2. **Integration Guide Modal**:
   - Click the `< /> Integration` button on any environment card.
   - Verify the modal opens with tabs for React, Node.js, Python, and Go.
   - Switch tabs and verify that the environment's key is pre-filled inside the code blocks.
   - Click "Copy Code Snippet" and verify the code is copied.
