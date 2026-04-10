using System.Windows;
using System.Windows.Controls;
using OrangeDesktop.ViewModels;

namespace OrangeDesktop.Views;

public partial class LoginWindow : Window
{
    private readonly LoginViewModel _viewModel;

    public LoginWindow(LoginViewModel viewModel)
    {
        InitializeComponent();
        _viewModel = viewModel;
        DataContext = _viewModel;
        _viewModel.LoginSucceeded += OnLoginSucceeded;
    }

    private void OnLoginSucceeded()
    {
        Dispatcher.Invoke(() =>
        {
            DialogResult = true;
            Close();
        });
    }

    private void PasswordBox_PasswordChanged(object sender, RoutedEventArgs e)
    {
        _viewModel.Password = ((PasswordBox)sender).Password;
    }
}
