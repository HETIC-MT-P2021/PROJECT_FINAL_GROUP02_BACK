package test



// type FakeSession struct { 
// 	status string 
// 	idle int 
// } 

// func (f *FakeSession) UpdateStatus(idle int, game string) error { 
// 	f.idle, f.status = idle, game 
// 	return nil 
// } 

// func TestStatusIsUpdated(t *testing.T) { 
// 	readyDependency := &discordgo.Ready{} 
// 	fakeSession := &FakeSession{} 

// 	ready(fakeSession, readyDependency) 

// 	// @todo assert that idle/game status were set to correct values 
// } 