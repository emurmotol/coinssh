package actions

func (as *ActionSuite) Test_Get_Home() {
	res := as.HTML("/").Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "Home")
}
