package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestBatchService_CreateJob(t *testing.T) {
	setup()
	defer teardown()
	uuid_str := uuid.New().String()
	opt := &BatchCreateJobOptions{
		ClientRequestToken: uuid_str,
		Description:        "test batch",
		Manifest: &BatchJobManifest{
			Location: &BatchJobManifestLocation{
				ETag:      "15150651828fa9cdcb8356b6d1c7638b",
				ObjectArn: "qcs::cos:ap-chengdu:uid/1250000000:sourcebucket-1250000000/manifests/batch-copy-manifest.csv",
			},
			Spec: &BatchJobManifestSpec{
				Fields: []string{"Bucket", "Key"},
				Format: "COSBatchOperations_CSV_V1",
			},
		},
		Operation: &BatchJobOperation{
			PutObjectCopy: &BatchJobOperationCopy{
				TargetResource: "qcs::cos:ap-chengdu:uid/1250000000:destinationbucket-1250000000",
			},
		},
		Priority: 1,
		Report: &BatchJobReport{
			Bucket:      "qcs::cos:ap-chengdu:uid/1250000000:sourcebucket-1250000000",
			Enabled:     "true",
			Format:      "Report_CSV_V1",
			Prefix:      "job-result",
			ReportScope: "AllTasks",
		},
		RoleArn: "qcs::cam::uin/100000000001:roleName/COS_Batch_QcsRole",
	}

	mux.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		testHeader(t, r, "x-cos-appid", "1250000000")
		testMethod(t, r, http.MethodPost)
		v := new(BatchCreateJobOptions)
		xml.NewDecoder(r.Body).Decode(v)

		want := opt
		want.XMLName = xml.Name{Local: "CreateJobRequest"}
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Batch.CreateJob request body: %+v, want %+v", v, want)
		}
		fmt.Fprint(w, `<?xml version='1.0' encoding='utf-8' ?>
<CreateJobResult>
    <JobId>53dc6228-c50b-46f7-8ad7-65e7159f1aae</JobId>
</CreateJobResult>`)
	})

	headers := &BatchRequestHeaders{
		XCosAppid: 1250000000,
	}
	ref, _, err := client.Batch.CreateJob(context.Background(), opt, headers)
	if err != nil {
		t.Fatalf("Batch.CreateJob returned error: %v", err)
	}

	want := &BatchCreateJobResult{
		XMLName: xml.Name{Local: "CreateJobResult"},
		JobId:   "53dc6228-c50b-46f7-8ad7-65e7159f1aae",
	}

	if !reflect.DeepEqual(ref, want) {
		t.Errorf("Batch.CreateJob returned %+v, want %+v", ref, want)
	}

}

func TestBatchService_DescribeJob(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/jobs/53dc6228-c50b-46f7-8ad7-65e7159f1aae", func(w http.ResponseWriter, r *http.Request) {
		testHeader(t, r, "x-cos-appid", "1250000000")
		testMethod(t, r, http.MethodGet)

		fmt.Fprint(w, `<?xml version='1.0' encoding='utf-8' ?>
<DescribeJobResult>
    <Job>
        <ConfirmationRequired>false</ConfirmationRequired>
        <CreationTime>2019-12-19T18:00:30Z</CreationTime>
        <Description>example-job</Description>
        <FailureReasons>
            <JobFailure>
                <FailureCode/>
                <FailureReason/>
            </JobFailure>
        </FailureReasons>
        <JobId>53dc6228-c50b-46f7-8ad7-65e7159f1aae</JobId>
        <Manifest>
            <Location>
                <ETag>&quot;15150651828fa9cdcb8356b6d1c7638b&quot;</ETag>
                <ObjectArn>qcs::cos:ap-chengdu:uid/1250000000:sourcebucket-1250000000/manifests/batch-copy-manifest.csv</ObjectArn>
            </Location>
            <Spec>
                <Fields>
                    <member>Bucket</member>
                    <member>Key</member>
                </Fields>
                <Format>COSBatchOperations_CSV_V1</Format>
            </Spec>
        </Manifest>
        <Operation>
            <COSPutObjectCopy>
                <TargetResource>qcs::cos:ap-chengdu:uid/1250000000:destinationbucket-1250000000</TargetResource>
            </COSPutObjectCopy>
        </Operation>
        <Priority>10</Priority>
        <ProgressSummary>
            <NumberOfTasksFailed>0</NumberOfTasksFailed>
            <NumberOfTasksSucceeded>10</NumberOfTasksSucceeded>
            <TotalNumberOfTasks>10</TotalNumberOfTasks>
        </ProgressSummary>
        <Report>
            <Bucket>qcs::cos:ap-chengdu:uid/1250000000:sourcebucket-1250000000</Bucket>
            <Enabled>true</Enabled>
            <Format>Report_CSV_V1</Format>
            <Prefix>job-result</Prefix>
            <ReportScope>AllTasks</ReportScope>
        </Report>
        <RoleArn>qcs::cam::uin/100000000001:roleName/COS_Batch_QcsRole</RoleArn>
        <Status>Complete</Status>
        <StatusUpdateReason>Job complete</StatusUpdateReason>
        <TerminationDate>2019-12-19T18:00:42Z</TerminationDate>
    </Job>
</DescribeJobResult>`)
	})

	headers := &BatchRequestHeaders{
		XCosAppid: 1250000000,
	}
	ref, _, err := client.Batch.DescribeJob(context.Background(), "53dc6228-c50b-46f7-8ad7-65e7159f1aae", headers)
	if err != nil {
		t.Fatalf("Batch.DescribeJob returned error: %v", err)
	}

	want := &BatchDescribeJobResult{
		XMLName: xml.Name{Local: "DescribeJobResult"},
		Job: &BatchDescribeJob{
			ConfirmationRequired: "false",
			CreationTime:         "2019-12-19T18:00:30Z",
			Description:          "example-job",
			FailureReasons:       &BatchJobFailureReasons{},
			JobId:                "53dc6228-c50b-46f7-8ad7-65e7159f1aae",
			Manifest: &BatchJobManifest{
				Location: &BatchJobManifestLocation{
					ETag:      "\"15150651828fa9cdcb8356b6d1c7638b\"",
					ObjectArn: "qcs::cos:ap-chengdu:uid/1250000000:sourcebucket-1250000000/manifests/batch-copy-manifest.csv",
				},
				Spec: &BatchJobManifestSpec{
					Fields: []string{"Bucket", "Key"},
					Format: "COSBatchOperations_CSV_V1",
				},
			},
			Operation: &BatchJobOperation{
				PutObjectCopy: &BatchJobOperationCopy{
					TargetResource: "qcs::cos:ap-chengdu:uid/1250000000:destinationbucket-1250000000",
				},
			},
			Priority: 10,
			ProgressSummary: &BatchProgressSummary{
				NumberOfTasksFailed:    0,
				NumberOfTasksSucceeded: 10,
				TotalNumberOfTasks:     10,
			},
			Report: &BatchJobReport{
				Bucket:      "qcs::cos:ap-chengdu:uid/1250000000:sourcebucket-1250000000",
				Enabled:     "true",
				Format:      "Report_CSV_V1",
				Prefix:      "job-result",
				ReportScope: "AllTasks",
			},
			RoleArn:            "qcs::cam::uin/100000000001:roleName/COS_Batch_QcsRole",
			Status:             "Complete",
			StatusUpdateReason: "Job complete",
			TerminationDate:    "2019-12-19T18:00:42Z",
		},
	}

	if !reflect.DeepEqual(ref, want) {
		t.Errorf("Batch.DescribeJob returned %+v, want %+v", ref, want)
	}
}

func TestBatchService_ListJobs(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		testHeader(t, r, "x-cos-appid", "1250000000")
		testMethod(t, r, http.MethodGet)
		vs := values{
			"maxResults": "2",
		}
		testFormValues(t, r, vs)

		fmt.Fprint(w, `<?xml version='1.0' encoding='utf-8' ?>
<ListJobsResult>
    <Jobs>
        <member>
            <CreationTime>2019-12-19T11:05:40Z</CreationTime>
            <Description>example-job</Description>
            <JobId>021140d8-67ca-4e89-8089-0de9a1e40943</JobId>
            <Operation>COSPutObjectCopy</Operation>
            <Priority>10</Priority>
            <ProgressSummary>
                <NumberOfTasksFailed>0</NumberOfTasksFailed>
                <NumberOfTasksSucceeded>10</NumberOfTasksSucceeded>
                <TotalNumberOfTasks>10</TotalNumberOfTasks>
            </ProgressSummary>
            <Status>Complete</Status>
            <TerminationDate>2019-12-19T11:05:56Z</TerminationDate>
        </member>
        <member>
            <CreationTime>2019-12-19T11:07:05Z</CreationTime>
            <Description>example-job</Description>
            <JobId>066d919e-49b9-429e-b844-e17ea7b16421</JobId>
            <Operation>COSPutObjectCopy</Operation>
            <Priority>10</Priority>
            <ProgressSummary>
                <NumberOfTasksFailed>0</NumberOfTasksFailed>
                <NumberOfTasksSucceeded>10</NumberOfTasksSucceeded>
                <TotalNumberOfTasks>10</TotalNumberOfTasks>
            </ProgressSummary>
            <Status>Complete</Status>
            <TerminationDate>2019-12-19T11:07:21Z</TerminationDate>
        </member>
    </Jobs>
    <NextToken>066d919e-49b9-429e-b844-e17ea7b16421</NextToken>
</ListJobsResult>`)
	})

	opt := &BatchListJobsOptions{
		MaxResults: 2,
	}
	headers := &BatchRequestHeaders{
		XCosAppid: 1250000000,
	}

	ref, _, err := client.Batch.ListJobs(context.Background(), opt, headers)
	if err != nil {
		t.Fatalf("Batch.DescribeJob returned error: %v", err)
	}

	want := &BatchListJobsResult{
		XMLName: xml.Name{Local: "ListJobsResult"},
		Jobs: &BatchListJobs{
			Members: []BatchListJobsMember{
				{
					CreationTime: "2019-12-19T11:05:40Z",
					Description:  "example-job",
					JobId:        "021140d8-67ca-4e89-8089-0de9a1e40943",
					Operation:    "COSPutObjectCopy",
					Priority:     10,
					ProgressSummary: &BatchProgressSummary{
						NumberOfTasksFailed:    0,
						NumberOfTasksSucceeded: 10,
						TotalNumberOfTasks:     10,
					},
					Status:          "Complete",
					TerminationDate: "2019-12-19T11:05:56Z",
				},
				{
					CreationTime: "2019-12-19T11:07:05Z",
					Description:  "example-job",
					JobId:        "066d919e-49b9-429e-b844-e17ea7b16421",
					Operation:    "COSPutObjectCopy",
					Priority:     10,
					ProgressSummary: &BatchProgressSummary{
						NumberOfTasksFailed:    0,
						NumberOfTasksSucceeded: 10,
						TotalNumberOfTasks:     10,
					},
					Status:          "Complete",
					TerminationDate: "2019-12-19T11:07:21Z",
				},
			},
		},
		NextToken: "066d919e-49b9-429e-b844-e17ea7b16421",
	}
	if !reflect.DeepEqual(ref, want) {
		t.Errorf("Batch.ListJobs returned %+v, want %+v", ref, want)
	}
}

func TestBatchService_UpdateJobsPriority(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/jobs/021140d8-67ca-4e89-8089-0de9a1e40943/priority", func(w http.ResponseWriter, r *http.Request) {
		testHeader(t, r, "x-cos-appid", "1250000000")
		testMethod(t, r, http.MethodPost)
		vs := values{
			"priority": "10",
		}
		testFormValues(t, r, vs)

		fmt.Fprint(w, `<?xml version='1.0' encoding='utf-8' ?>
<UpdateJobPriorityResult>
    <JobId>021140d8-67ca-4e89-8089-0de9a1e40943</JobId>
    <Priority>10</Priority>
</UpdateJobPriorityResult>`)
	})

	opt := &BatchUpdatePriorityOptions{
		JobId:    "021140d8-67ca-4e89-8089-0de9a1e40943",
		Priority: 10,
	}

	headers := &BatchRequestHeaders{
		XCosAppid: 1250000000,
	}

	ref, _, err := client.Batch.UpdateJobPriority(context.Background(), opt, headers)
	if err != nil {
		t.Fatalf("Batch.UpdateJobPriority returned error: %v", err)
	}

	want := &BatchUpdatePriorityResult{
		XMLName:  xml.Name{Local: "UpdateJobPriorityResult"},
		JobId:    "021140d8-67ca-4e89-8089-0de9a1e40943",
		Priority: 10,
	}
	if !reflect.DeepEqual(ref, want) {
		t.Errorf("Batch.UpdateJobsPriority returned %+v, want %+v", ref, want)
	}
}

func TestBatchService_UpdateJobsStatus(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/jobs/021140d8-67ca-4e89-8089-0de9a1e40943/status", func(w http.ResponseWriter, r *http.Request) {
		testHeader(t, r, "x-cos-appid", "1250000000")
		testMethod(t, r, http.MethodPost)
		vs := values{
			"requestedJobStatus": "Ready",
			"statusUpdateReason": "to do",
		}
		testFormValues(t, r, vs)

		fmt.Fprint(w, `<?xml version='1.0' encoding='utf-8' ?>
<UpdateJobStatusResult>
    <JobId>021140d8-67ca-4e89-8089-0de9a1e40943</JobId>
    <Status>Ready</Status>
    <StatusUpdateReason>to do</StatusUpdateReason>
</UpdateJobStatusResult>`)
	})

	opt := &BatchUpdateStatusOptions{
		JobId:              "021140d8-67ca-4e89-8089-0de9a1e40943",
		RequestedJobStatus: "Ready",
		StatusUpdateReason: "to do",
	}

	headers := &BatchRequestHeaders{
		XCosAppid: 1250000000,
	}

	ref, _, err := client.Batch.UpdateJobStatus(context.Background(), opt, headers)
	if err != nil {
		t.Fatalf("Batch.UpdateJobStatus returned error: %v", err)
	}

	want := &BatchUpdateStatusResult{
		XMLName:            xml.Name{Local: "UpdateJobStatusResult"},
		JobId:              "021140d8-67ca-4e89-8089-0de9a1e40943",
		Status:             "Ready",
		StatusUpdateReason: "to do",
	}
	if !reflect.DeepEqual(ref, want) {
		t.Errorf("Batch.UpdateJobsStatus returned %+v, want %+v", ref, want)
	}
}

func TestBatchService_DeleteJob(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/jobs/53dc6228-c50b-46f7-8ad7-65e7159f1aae", func(w http.ResponseWriter, r *http.Request) {
		testHeader(t, r, "x-cos-appid", "1250000000")
		testMethod(t, r, http.MethodDelete)
	})

	headers := &BatchRequestHeaders{
		XCosAppid: 1250000000,
	}
	_, err := client.Batch.DeleteJob(context.Background(), "53dc6228-c50b-46f7-8ad7-65e7159f1aae", headers)
	if err != nil {
		t.Fatalf("Batch.DescribeJob returned error: %v", err)
	}
}
