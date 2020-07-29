package versions_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.ibm.com/coligo/satcon-client/client/actions"
	. "github.ibm.com/coligo/satcon-client/client/actions/versions"
	"github.ibm.com/coligo/satcon-client/client/web/webfakes"
)

var _ = Describe("AddChannelVersion", func() {

	var (
		orgID, channelUuid, name string
		description, token       string
		content                  []byte
		c                        VersionService
		h                        *webfakes.FakeHTTPClient
		response                 *http.Response
	)

	BeforeEach(func() {
		orgID = "someorg"
		channelUuid = "somechannel"
		name = "somechannel"
		content = []byte("YXBpVmVyc2lvbjogdjEKa2luZDogUG9kCm1ldGFkYXRhOgogIG5hbWU6IHZlcnNpb250ZXN0CnNwZWM6CiAgY29udGFpbmVyczoKICAtIG5hbWU6IHZlcnNpb250ZXN0CiAgICBpbWFnZTogaHR0cGQ6YWxwaW5lCg==")
		description = "somedescription"
		token = "thisissupposedtobeatoken"

		h = &webfakes.FakeHTTPClient{}
		response = &http.Response{}
		h.DoReturns(response, nil)
	})

	JustBeforeEach(func() {
		c, _ = NewClient("https://foo.bar", h)
		Expect(c).NotTo(BeNil())

		Expect(h.DoCallCount()).To(Equal(0))
	})

	Describe("NewAddChannelVersionVariables", func() {
		It("Returns a correctly configured set of variables", func() {
			vars := NewAddChannelVersionVariables(orgID, channelUuid, name, ContentType, string(content), "", description)
			Expect(vars.Type).To(Equal(actions.QueryTypeMutation))
			Expect(vars.QueryName).To(Equal(QueryAddChannelVersion))
			Expect(vars.OrgID).To(Equal(orgID))
			Expect(vars.Name).To(Equal(name))
			Expect(vars.Args).To(Equal(map[string]string{
				"orgId":       "String!",
				"channelUuid": "String!",
				"name":        "String!",
				"type":        "String!",
				"content":     "String",
				"file":        "Upload",
				"description": "String",
			}))
			Expect(vars.Returns).To(ConsistOf(
				"versionUuid",
				"success",
			))
		})
	})

	Describe("AddChannelVersion", func() {

		var (
			addChannelVersionResponse AddChannelVersionResponse
		)

		BeforeEach(func() {
			addChannelVersionResponse = AddChannelVersionResponse{
				Data: &AddChannelVersionResponseData{
					Details: &AddChannelVersionResponseDataDetails{
						VersionUUID: "newversionuuid",
						Success:     true,
					},
				},
			}

			respBodyBytes, err := json.Marshal(addChannelVersionResponse)
			Expect(err).NotTo(HaveOccurred())

			response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
		})

		It("Sends the http request", func() {
			_, err := c.AddChannelVersion(orgID, channelUuid, name, content, description, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(h.DoCallCount()).To(Equal(1))
		})

		It("Returns the add channel version details", func() {
			details, _ := c.AddChannelVersion(orgID, channelUuid, name, content, description, token)
			Expect(details).NotTo(BeNil())

			expected := addChannelVersionResponse.Data.Details
			Expect(*details).To(Equal(*expected))
		})

	})

})