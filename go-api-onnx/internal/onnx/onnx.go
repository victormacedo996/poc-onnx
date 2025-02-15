package onnx

import ort "github.com/yalue/onnxruntime_go"

var MAPPED_LABELS = map[int64]string{
	0: "iris-setosa",
	1: "iris-versicolor",
	2: "iris-virginica",
}

func InitializeEnvironment(shared_lib_path string) {
	ort.SetSharedLibraryPath(shared_lib_path)
	err := ort.InitializeEnvironment()
	if err != nil {
		panic(err)
	}
}

func CreateSession(model_path string) (*ort.DynamicAdvancedSession, error) {
	options := ort.SessionOptions{}
	defer options.Destroy()
	sess, err := ort.NewDynamicAdvancedSession(model_path, []string{"float_input"}, []string{"output_label", "output_probability"}, &options)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func Create_input_tensor(input_values []float32, dimensions ...int64) (*ort.Tensor[float32], error) {
	input_shape := ort.NewShape(dimensions...)
	input_tensor, err := ort.NewTensor(input_shape, input_values)
	if err != nil {
		return &ort.Tensor[float32]{}, err
	}
	return input_tensor, nil
}

func Predict(input_values []float32, sess *ort.DynamicAdvancedSession, input_shape ...int64) (string, error) {
	input_tensor, err := Create_input_tensor(input_values, input_shape...)
	if err != nil {
		return "", err
	}
	defer input_tensor.Destroy()

	output_values := []ort.Value{nil, nil}

	err = sess.Run([]ort.Value{input_tensor}, output_values)
	if err != nil {
		return "", err
	}
	defer output_values[0].Destroy()
	defer output_values[1].Destroy()

	label_tensor := output_values[0].(*ort.Tensor[int64])
	defer label_tensor.Destroy()

	return MAPPED_LABELS[label_tensor.GetData()[0]], nil

}

func Predict_proba(input_values []float32, sess *ort.DynamicAdvancedSession, input_shape ...int64) (map[string]float32, error) {
	input_tensor, err := Create_input_tensor(input_values, input_shape...)
	if err != nil {
		return map[string]float32{}, err
	}
	defer input_tensor.Destroy()

	output_values := []ort.Value{nil, nil}

	err = sess.Run([]ort.Value{input_tensor}, output_values)
	if err != nil {
		return map[string]float32{}, err
	}
	defer output_values[0].Destroy()
	defer output_values[1].Destroy()

	sequence := output_values[1].(*ort.Sequence)
	probability_maps, err := sequence.GetValues()
	if err != nil {
		panic(err)
	}

	probabilities := make(map[string]float32)

	for i := range probability_maps {
		m := probability_maps[i].(*ort.Map)
		keys, values, err := m.GetKeysAndValues()
		if err != nil {
			panic(err)
		}
		keys_tensor := keys.(*ort.Tensor[int64])
		values_tensor := values.(*ort.Tensor[float32])

		for j, key := range keys_tensor.GetData() {
			value := values_tensor.GetData()[j]
			probabilities[MAPPED_LABELS[key]] = value
		}
	}

	return probabilities, nil
}
