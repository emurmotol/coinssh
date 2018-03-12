package actions

func (as *ActionSuite) Test_Get_Routes() {
	res := as.HTML("/routes").Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "Welcome to Buffalo")
}