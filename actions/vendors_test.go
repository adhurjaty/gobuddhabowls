package actions

func (as *ActionSuite) Test_ListVendor_View() {
	vendor, err := createVendor(as.DB)
	as.NoError(err)

	res := as.HTML("/vendors").Get()
	as.Contains(res.Body.String(), vendor.ID.String())
}

func (as *ActionSuite) Test_VendorsResource_Show() {
	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_VendorsResource_New() {
	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_VendorsResource_Create() {
	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_VendorsResource_Edit() {
	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_VendorsResource_Update() {
	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_VendorsResource_Destroy() {
	as.Fail("Not Implemented!")
}
