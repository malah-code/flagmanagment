# Quickstart Validation Guide

## Prerequisites
- Backend running on `localhost:8080`
- Frontend running on `localhost:3000`
- `MailHog` running on `localhost:8025` for local email testing (can be spun up via Docker).

## Validation Steps

### 1. Configure SMTP
- Navigate to the UI `http://localhost:3000/settings/system` (new route to be built).
- Enter `localhost` for host, `1025` for port (MailHog's SMTP port).
- Click **Test Connection**. 
- Verify a success toast appears and check MailHog at `http://localhost:8025` for the test email.

### 2. Invite User
- Navigate to `http://localhost:3000/settings/users`.
- Click **Invite User**, enter `test@example.com`, select "Project Editor", and choose a project.
- Submit the form. 
- Verify the user appears as "Pending" in the list.

### 3. Verify Email
- Open MailHog at `http://localhost:8025`.
- Verify the invitation email was received containing a secure registration link.

### 4. Edit User
- Click the pencil edit icon next to the pending user.
- Change their role to "Global Administrator".
- Save changes and verify the table updates to reflect "All Projects" access.
