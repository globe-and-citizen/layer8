## Test Cases Checklist

### SECTION 1: RESOURCE SERVER (http://localhost:5001/)
- [ ] Created a new user on the Layer8 Portal

### SECTION 2: WE'VE GOT POEMS (http://localhost:5173/)
- [ ] Login as 'tester', '1234'
- [ ] Choose login with layer8 
- [ ] Logging in with the newly registered credentials from Section 1 succeeds in popup
- [ ] User chooses to share their new "Username" & "Country" from the Layer8 Resource Server
- [ ] Choose poems 1 - 3 works correctly 

### SECTION 3: IMSHARER (http://localhost:5174/)
- [ ] Uploads an image works
- [ ] Reload leads to instant reload (demonstrating proper caching)


